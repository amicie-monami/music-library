package app

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/amicie-monami/music-library/config"
	"github.com/amicie-monami/music-library/internal/server"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func Run(ctx context.Context, config *config.Config) {
	db := databaseConnect(config.Database.Source)
	slog.Info("successful connection to the database")

	runMigrations(db.DB)

	server := server.New(ctx, config, db)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		slog.Info("starting server", "addr", config.Server.Addr)
		if err := server.Run(ctx); err != nil {
			slog.Error("http server", "msg", err)
		}
	}()

	//checks the context for termination signals
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("http server shutdown", "msg", err)
		}
		slog.Info("http server was successfully shutdown")
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

// runMigrations starts the database migration proccess
func runMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		version, dirty, _ := m.Version()
		log.Fatalf("failed to apply migrations, version=%d, dirty=%t err=%s", version, dirty, err)
	}

	version, _, _ := m.Version()
	slog.Info("apply migrations", "version", version)
}
