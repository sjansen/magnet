package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	cmw "github.com/go-chi/chi/middleware"

	"github.com/sjansen/magnet/internal/aws"
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
	)
	if s.useSCS {
		r.Use(s.sessions.LoadAndSave)
		r.Use(s.relaystate.LoadAndSave)
	}

	requireLogin := s.saml.RequireAccount
	s3 := aws.NewS3Client(s.aws)

	r.Get("/", handlers.Root)
	r.Mount("/saml/", s.saml)
	r.Handle("/browse/*", requireLogin(
		handlers.NewBrowser("/browse/", s.config, s3),
	))
	r.Handle("/upload/", requireLogin(
		handlers.NewUploader("/upload/", s.config.Bucket, s3),
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
