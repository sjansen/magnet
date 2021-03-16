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
	r.Get("/", handlers.Root)
	r.Mount("/saml/", s.saml)
	r.Handle("/browse/*", requireLogin(
		handlers.NewBrowser("/browse/", s.config, s.config.AWS.NewS3Client()),
	))
	r.Handle("/upload/", requireLogin(
		handlers.NewUploader("/upload/", s.config.Bucket, s.config.AWS.Config),
	))
	r.Handle("/whoami", requireLogin(
		http.HandlerFunc(handlers.WhoAmI),
	))

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
