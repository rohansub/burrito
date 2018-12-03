package db

import (
	"github.com/go-redis/redis"
	"github.com/rcsubra2/burrito/src/utils"
	"time"
)

// Database - The interface for database clients
type Database interface {
	Get(args []string) map[string]string
	Set(items []utils.Pair) bool
	Delete(args []string) bool
}

// RedisDBInterface - interface that is implemented by redis.Client,
// and the RedisMockClient
type RedisDBInterface interface {
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(keys ...string) *redis.IntCmd
}
