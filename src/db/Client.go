package db

import (
	"github.com/go-redis/redis"
	"github.com/rcsubra2/burrito/src/environment"
)

// DBClient -
type DBClient interface {
	Get(GetReq, []*environment.Env) map[string]string
}



type RedisClient struct {
	db *redis.Client
}

func NewRedisClient(uri string) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// Output: PONG <nil>
	return &RedisClient {
		db: client,
	}

}


func (rc *RedisClient) Get(req GetReq, envs []*environment.Env) map[string]string {
	keys := make([]string, len(req.ArgNames))
	for i, ar := range req.ArgNames {
		val, ok := ar.GetValue(envs)
		if ok {
			keys[i] = val
		}
	}

	respData := make(map[string]string)

	for _, k := range keys {
		val, err := rc.db.Get(k).Result()
		if err == nil {
			respData[k] = val
		}

	}

	return respData

}
