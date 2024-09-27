package server

import (
	"context"
	"net/http"
	"time"

	"github.com/amicie-monami/music-library/config"
	"github.com/amicie-monami/music-library/internal/handler/v1"
	"github.com/amicie-monami/music-library/pkg/router"
)

// server wraps the http.Server
type server struct {
	srv *http.Server
}

func New(config *config.Config) *server {
	router := router.New()
	configureRouter(router)
	srv := &http.Server{
		Addr:           config.Server.Addr,
		ReadTimeout:    time.Duration(config.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.Server.WriteTimeout) * time.Second,
		IdleTimeout:    time.Duration(config.Server.IdleTimeout) * time.Second,
		MaxHeaderBytes: config.Server.MaxHeaderBytes,
		Handler:        router,
	}
	return &server{srv}
}

func configureRouter(router *router.Router) {
	router.Handle("api/v1/songs", handler.AddSong()).Methods("POST")
}

// Run starts the server
func (s *server) Run(ctx context.Context) error {
	return s.srv.ListenAndServe()
}

// Shutdown performs a graceful shutdown of the server
func (s *server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
