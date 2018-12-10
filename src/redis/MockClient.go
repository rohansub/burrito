package redis

import (
	"github.com/go-redis/redis"
	"errors"
	"time"
)

type Client struct {
	data map[string]string
}

// NewMockRedisClient - Create Mock Redis client
func NewMockRedisClient(init map[string]string) *Client{
	return &Client {
		data: init,
	}

}

func (c * Client) Get(key string) *redis.StringCmd {
	val, ok := c.data[key]
	if !ok {
		return redis.NewStringResult("", errors.New("Key not found"))
	} else {
		return redis.NewStringResult(val, nil)
	}
}


func (c * Client) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	c.data[key] = value.(string)
	return redis.NewStatusCmd()
}

func (c* Client) Del(keys ...string) *redis.IntCmd {
	val := 0
	for _, k := range keys {
		_, ok := c.data[k]
		if ok {
			delete(c.data, k)
			val += 1
		}
	}
	return redis.NewIntCmd(val)
}
