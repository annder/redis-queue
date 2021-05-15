package client

import "github.com/go-redis/redis/v8"

func Reids() *redis.Client {
	return redis.NewClient(&redis.Options{
		Password: "",
		Addr:     "localhost:6379",
		DB:       0,
	})
}
