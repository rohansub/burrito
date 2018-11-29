package db

import (
	"github.com/go-redis/redis"
	"github.com/rcsubra2/burrito/src/environment"
)

// Database - The interface for database clients
type Database interface {
	Get(req GetReq, envs []*environment.Env) map[string]string
}

// RedisDBInterface - interface that is implemented by redis.Client,
// and the RedisMockClient
type RedisDBInterface interface {
	Get(key string) *redis.StringCmd
}
