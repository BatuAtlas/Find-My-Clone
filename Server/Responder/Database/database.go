package database

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var once sync.Once
var pool *pgxpool.Pool //an address

func Initialize() {
	once.Do(func() {
		dbpool, err := pgxpool.New(context.Background(), "postgres://postgres:newpassword@localhost:5432/FindMy")
		if err != nil {
			log.Fatal(err)
		}

		pool = dbpool
	})
}

func GetPool() *pgxpool.Pool {
	return pool
}
