package db

import (
	"github.com/go-redis/redis"
)

// Database - The interface for database clients
type Database interface {
	Get(args []string) map[string]string
}

// RedisDBInterface - interface that is implemented by redis.Client,
// and the RedisMockClient
type RedisDBInterface interface {
	Get(key string) *redis.StringCmd
}
