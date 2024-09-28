package app

import (
	"context"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/amicie-monami/music-library/config"
	"github.com/amicie-monami/music-library/internal/repo"
	"github.com/amicie-monami/music-library/internal/server"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func Run(ctx context.Context, config *config.Config) {
	db := databaseConnect(config.Database.Source)
	songRepo := repo.NewSong(db)
	server := server.New(config, songRepo)

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
		slog.Debug("http server was successfully shutdown")
	}()

	wg.Wait()
}

// databaseConnect connects to the database at the source address and pings it
func databaseConnect(source string) *sqlx.DB {
	db, err := sqlx.Open("pgx", source)
	if err != nil || db.Ping() != nil {
		log.Fatalf("failed to connect to database, msg=%s", err)
	}
	return db
}
