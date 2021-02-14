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
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/crewjam/saml/samlsp"
	"github.com/go-chi/chi"

	"github.com/sjansen/magnet/internal/config"
)

var _ samlsp.RequestTracker = &Server{}
var _ samlsp.Session = &Server{}

// Server provides Strongbox's API
type Server struct {
	config     *config.Config
	lambda     *chiadapter.ChiLambda
	relaystate *scs.SessionManager
	router     *chi.Mux
	saml       *samlsp.Middleware
	sessions   *scs.SessionManager

	useSCS bool

	done chan struct{}
	wg   sync.WaitGroup
}

// New creates a new Server
func New(cfg *config.Config) (*Server, error) {
	s := &Server{
		config: cfg,
		done:   make(chan struct{}),
	}

	sp, err := newSAMLMiddleware(cfg)
	if err != nil {
		return nil, err
	}
	s.saml = sp

	relaystate, sessions, err := s.openDynamoStores(cfg)
	if err != nil {
		return nil, err
	}
	s.addSCS(relaystate, sessions)

	s.addRouter()
	return s, nil
}

// HandleLambda starts the server waiting for events from AWS Lambda.
func (s *Server) HandleLambda() {
	s.lambda = chiadapter.New(s.router)
	lambda.Start(s.LambdaHandler)
}

// LambdaHandler processes a single Lambda event.
func (s *Server) LambdaHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return s.lambda.ProxyWithContext(ctx, req)
}

// ListenAndServe starts the server waiting for network connections.
func (s *Server) ListenAndServe() error {
	fmt.Println("Listening to", s.config.Root.String())

	server := &http.Server{
		Addr:    s.config.Addr,
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
