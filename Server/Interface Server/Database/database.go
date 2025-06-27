package database

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"mime/multipart"
	"regexp"
	"sync"

	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

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

func RemoveValue(array []int64, remove int64) []int64 { //chatgpt function
	result := array[:0]
	for _, v := range array {
		if v != remove {
			result = append(result, v)
		}
	}
	return result
}

func IsValidEmail(email string) bool {
	// Basit ve pratik bir e-posta regex’i
	var re = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func IsSHA256Hash(s string) bool {
	matched, _ := regexp.MatchString("^[a-f0-9]{64}$", s)
	return matched
}

func GenerateRandomFilename() (string, error) {
	b := make([]byte, 25)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func IsValidImage(fileHeader *multipart.FileHeader) bool {
	f, err := fileHeader.Open()
	if err != nil {
		return false
	}
	defer f.Close()

	_, _, err = image.Decode(f)
	return err == nil
}
