package util

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var redisCon *redis.Client = nil
var ctx = context.Background()

func connect() *redis.Client {
	if redisCon == nil {
		redisCon = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
	}
	return redisCon
}

func RdbSet(key string, value string) {
	r := connect()
	err := r.Set(ctx, key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}

func RdbGet(key string) string {
	r := connect()
	val, err := r.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return val
}

func RdbWatch(key string) {
	r := connect()
	r.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		return nil
	})
}
