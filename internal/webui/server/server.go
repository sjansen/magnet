package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/alexedwards/scs/v2"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	proxy "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/crewjam/saml/samlsp"
	"github.com/go-chi/chi"

	"github.com/sjansen/magnet/internal/config"
)

var _ samlsp.RequestTracker = &Server{}
var _ samlsp.Session = &Server{}

// Server provides Magnet's API
type Server struct {
	config     *config.WebUI
	lambda     *proxy.ChiLambda
	relaystate *scs.SessionManager
	router     *chi.Mux
	saml       *samlsp.Middleware
	sessions   *scs.SessionManager

	done chan struct{}
	wg   sync.WaitGroup
}

// New creates a new Server
func New(ctx context.Context) (*Server, error) {
	fmt.Println("Loading config...")
	cfg, err := config.LoadWebUIConfig(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Println("Loading SAML config...")
	sp, err := newSAMLMiddleware(cfg)
	if err != nil {
		return nil, err
	}

	s := &Server{
		config: cfg,
		done:   make(chan struct{}),
		saml:   sp,
	}

	fmt.Println("Preparing session store...")
	relaystate, sessions, err := s.openDynamoStores(cfg)
	if err != nil {
		return nil, err
	}
	s.addSCS(relaystate, sessions)

	fmt.Println("Configuring routes...")
	s.addRouter()
	return s, nil
}

// LambdaHandler processes a single Lambda event.
func (s *Server) LambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return s.lambda.ProxyWithContext(ctx, req)
}

// ListenAndServe starts the server waiting for network connections.
func (s *Server) ListenAndServe() error {
	fmt.Printf("Listening to http://%s/\n", s.config.Listen)

	server := &http.Server{
		Addr:    s.config.Listen,
		Handler: s.router,
	}
	go func() {
		ch := make(chan os.Signal, 10)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		<-ch
		fmt.Println("Exiting...")
		_ = server.Shutdown(context.Background())
	}()

	err := server.ListenAndServe()
	close(s.done)
	s.wg.Wait()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// StartLambdaHandler starts the server waiting for events from AWS Lambda.
func (s *Server) StartLambdaHandler() {
	s.lambda = proxy.New(s.router)
	lambda.Start(s.LambdaHandler)
}
