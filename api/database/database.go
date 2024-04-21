package database

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func NewClient(dbNo int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDR"),
		Password: os.Getenv("DB_PASSWORD"),
		DB:       dbNo,
	})
}
