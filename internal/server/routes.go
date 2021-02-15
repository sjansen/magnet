package server

import (
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-chi/chi"
	cmw "github.com/go-chi/chi/middleware"

	"github.com/sjansen/magnet/internal/handlers"
)

func (s *Server) addRouter() {
	r := chi.NewRouter()
	s.router = r

	r.Use(
		cmw.RequestID,
		cmw.RealIP,
		cmw.Logger,
		cmw.Recoverer,
		cmw.Timeout(5*time.Second),
		cmw.Heartbeat("/ping"),
	)
	if s.useSCS {
		r.Use(s.sessions.LoadAndSave)
		r.Use(s.relaystate.LoadAndSave)
	}

	r.Get("/", handlers.Root)
	r.Mount("/saml/", s.saml)
	r.Handle("/browse/*", s.saml.RequireAccount(http.HandlerFunc(
		handlers.NewBrowser("/browse/", s.config.Bucket, s3.New(s.aws)).Handler,
	)))
	r.Handle("/profile", s.saml.RequireAccount(http.HandlerFunc(handlers.Profile)))
}
