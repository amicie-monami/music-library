package server

import (
	"context"
	"net/http"
	"time"

	"github.com/amicie-monami/music-library/config"
	"github.com/amicie-monami/music-library/internal/repo"
	"github.com/amicie-monami/music-library/internal/usecase"
	"github.com/gorilla/mux"
)

// server wraps the http.Server
type server struct {
	srv *http.Server
}

func New(config *config.Config, songRepo *repo.Song) *server {
	router := mux.NewRouter()
	configureRouter(router, songRepo, usecase.New(songRepo))
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

// Run starts the server
func (s *server) Run(ctx context.Context) error {
	return s.srv.ListenAndServe()
}

// Shutdown performs a graceful shutdown of the server
func (s *server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
