package db

import (
	"github.com/go-redis/redis"
	"github.com/m4tthewde/tmihooks/internal/structs"
)

func AddWebhook(webhook *structs.Webhook) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	for _, channel := range webhook.Channels {
		rdb.RPush(channel, webhook.URI)
	}
}

func GetURIs(channel string) ([]string, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	result, err := rdb.LRange(channel, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func Clear() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	rdb.FlushAll()
}
