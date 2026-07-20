package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"resuming/database/sqlc"
	systemconfig "resuming/system-config"
)

var Pool *pgxpool.Pool
var Queries *sqlc.Queries

func DatabaseConnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, systemconfig.DatabaseDSN)
	if err != nil {
		log.Println("Failed to create connection pool")
		log.Println(err)
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Println("Failed to ping database")
		log.Println(err)
		return fmt.Errorf("failed to ping database: %w", err)
	}

	Pool = pool
	Queries = sqlc.New(pool)

	log.Println("Successfully Connected to Database")
	return nil
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
