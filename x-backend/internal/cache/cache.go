package cache

import (
	"github.com/redis/go-redis/v9"

	"github.com/jourloy/X-Backend/internal/config/env"
)

var Client *redis.Client

// InitCache подключается к кэшу
func InitCache() {
	Client = redis.NewClient(&redis.Options{
		Addr:     env.CacheDSN,
		Password: ``,
		DB:       0,
	})
}
