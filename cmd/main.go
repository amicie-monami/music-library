package main

import (
	"context"
	"flag"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/amicie-monami/music-library/config"
	"github.com/amicie-monami/music-library/internal/app"

	_ "github.com/amicie-monami/music-library/docs"
)

// @title Music Library API
// @version 1.0
// @description API for managing songs in the music library.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	flag.Parse()
	config := config.MustLoadFromEnv()
	// subcribe on terminate and quit signals
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-ctx.Done()
		cancel()
	}()
	app.Run(ctx, config)
}
