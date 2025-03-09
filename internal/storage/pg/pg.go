package pg

import (
	"context"
	"errors"
	"fin-manager/internal/config"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Драйвер для файловых миграций
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type DB struct {
	Pool   *pgxpool.Pool
	Config *config.Config
}

func New(ctx context.Context, config config.Config) *DB {
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.Storage.DbUser,
		config.Storage.DbPassword,
		config.Storage.DbHost,
		config.Storage.DbPort,
		config.Storage.DbName)
	pool, err := pgxpool.New(ctx, databaseUrl)
	if err != nil {
		log.Fatalf("error connecting to postgres: %v", err)
	}
	db := &DB{
		Pool:   pool,
		Config: &config,
	}
	db.Ping(ctx)
	return db
}

func (db *DB) Close() {
	log.Printf("closing database connection")
	db.Pool.Close()
}

func (db *DB) Ping(ctx context.Context) {
	if err := db.Pool.Ping(ctx); err != nil {
		log.Fatalf("error pinging database: %v", err)
	}
	log.Printf("pinged database")
}

func (db *DB) Migrate(migrationPath string) error {
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", db.Config.Storage.DbUser, db.Config.Storage.DbPassword, db.Config.Storage.DbHost, db.Config.Storage.DbPort, db.Config.Storage.DbName)
	m, err := migrate.New(migrationPath, databaseUrl)
	fmt.Println(databaseUrl)
	if err != nil {
		return fmt.Errorf("error create the migrate instance: %v", err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("error migrate up: %v", err)
	}

	log.Println("migration done")

	return nil
}
