package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	cmw "github.com/go-chi/chi/middleware"

	"github.com/sjansen/magnet/internal/webui/handlers"
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
		s.sessions.LoadAndSave,
		s.relaystate.LoadAndSave,
	)

	requireLogin := s.saml.RequireAccount
	r.Method("GET", "/",
		handlers.NewRoot(s.config),
	)
	r.Method("GET", "/browse/*", requireLogin(
		handlers.NewBrowser("/browse/", s.config, s.config.AWS.NewS3Client()),
	))
	r.Method("GET", "/upload/", requireLogin(
		handlers.NewUploader("/upload/", s.config, s.config.Bucket),
	))
	r.Method("GET", "/whoami/", requireLogin(
		handlers.NewWhoAmI(s.config),
	))
	r.Mount("/saml/", s.saml)

	if s.config.Development {
		r.Get("/favicon.ico",
			http.StripPrefix("/",
				http.FileServer(http.Dir("terraform/modules/app/icons/")),
			).ServeHTTP,
		)
		r.Get("/magnet/icons/*",
			http.StripPrefix("/magnet/icons/",
				http.FileServer(http.Dir("terraform/modules/app/icons/")),
			).ServeHTTP,
		)
	}
}
