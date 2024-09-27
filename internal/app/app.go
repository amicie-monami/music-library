package app

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/amicie-monami/music-library/config"
	"github.com/amicie-monami/music-library/internal/server"
)

func Run(ctx context.Context, config *config.Config) {
	server := server.New(config)
	var wg sync.WaitGroup
	wg.Add(1)

	//http server startup
	go func() {
		defer wg.Done()
		slog.Info("starting server on", "addr", config.Server.Addr)
		if err := server.Run(ctx); err != nil {
			slog.Error("http server", "msg", err)
		}
	}()

	//checks the context for caught signals about terms
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("http server shutdown", "msg", err)
		}
	}()

	wg.Wait()
}
