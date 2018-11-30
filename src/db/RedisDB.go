package db

import (
	"github.com/rcsubra2/burrito/src/environment"
)

// RedisDB - A client that interacts with Redis data
type RedisDB struct {
	db   RedisDBInterface
}

// NewRedisDB - create RedisDB client, given uri
func NewRedisDB(client RedisDBInterface) *RedisDB {


	return &RedisDB{
		db: client,
	}

}

// Get - Perform GetReq req given a list of environments
func (rc *RedisDB) Get(req GetReq, envs []*environment.Env) map[string]string {
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
