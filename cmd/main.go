package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/amicie-monami/music-library/config"
	"github.com/amicie-monami/music-library/internal/app"
)

func main() {
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
