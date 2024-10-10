package database

import (
    "context"
    "github.com/jackc/pgx/v4/pgxpool"
    "log"
)

func InitDB() *pgxpool.Pool {
	dbUrl := "postgres://sa:root@localhost:5432/test"
    dbPool, err := pgxpool.Connect(context.Background(), dbUrl)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }
    return dbPool
}
