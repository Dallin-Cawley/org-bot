package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
)

// BeginTransaction begins a transaction from a connection of the application pgxpool.Pool.
func BeginTransaction(ctx context.Context) (pgx.Tx, error) {
	return pool.BeginTx(ctx, pgx.TxOptions{})
}

// Connect attempts to connect to the database designated by the dbURL. A total of 10 attempts will be made before
// an error is returned. The database module will be initialized after this method completes successfully
func Connect(ctx context.Context, dbURL string) error {
	maxAttempts := 10

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return err
	}

	pgxPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return err
	}

	for i := 0; i < maxAttempts; i++ {
		log.Printf("Waiting for database. Attempt %d", i)

		if attemptConnection(ctx, pgxPool) {
			pool = pgxPool
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	pgxPool.Close()
	return fmt.Errorf("unable to confirm database alive after %d seconds", maxAttempts)
}

// attemptConnection attempts to ping the database to confirm the database is up and running.
func attemptConnection(ctx context.Context, pool *pgxpool.Pool) bool {
	if err := pool.Ping(ctx); err != nil {
		fmt.Println("error connecting to database: " + err.Error())
		return false
	}

	return true
}
