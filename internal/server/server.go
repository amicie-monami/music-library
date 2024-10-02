package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/amicie-monami/music-library/config"
	"github.com/amicie-monami/music-library/internal/repository"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// server wraps the http.Server
type server struct {
	srv *http.Server
}

func New(ctx context.Context, config *config.Config, db *sqlx.DB) *server {
	router := mux.NewRouter()
	songRepo := repository.NewSong(db)

	configureRouter(router, songRepo)
	srv := &http.Server{
		Addr:           config.Server.Addr,
		ReadTimeout:    time.Duration(config.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(config.Server.WriteTimeout) * time.Second,
		IdleTimeout:    time.Duration(config.Server.IdleTimeout) * time.Second,
		MaxHeaderBytes: config.Server.MaxHeaderBytes,
		Handler:        router,
		BaseContext:    func(l net.Listener) context.Context { return ctx },
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
