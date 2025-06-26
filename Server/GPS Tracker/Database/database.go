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

func IsSubset(sub, full []int64) bool {
	lookup := make(map[int64]bool)
	for _, v := range full {
		lookup[v] = true
	}

	for _, v := range sub {
		if !lookup[v] {
			return false
		}
	}

	return true
}

func RemoveValue(s []int64, x int64) []int64 { //chatgpt function
	result := s[:0]
	for _, v := range s {
		if v != x {
			result = append(result, v)
		}
	}
	return result
}
