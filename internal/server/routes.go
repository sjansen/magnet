package server

import (
	"net/http"
	"time"

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
	r.Get("/public", handlers.Public)
	r.Mount("/saml/", s.saml)
	r.Handle("/profile", s.saml.RequireAccount(http.HandlerFunc(handlers.Profile)))
	r.Handle("/secret", s.saml.RequireAccount(http.HandlerFunc(handlers.Secret)))
}
