package db

import (
	"github.com/rcsubra2/burrito/src/utils"
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

// Get - Perform get on redis database given a list of strings
func (rc *RedisDB) Get(keys []string) map[string]string {
	respData := make(map[string]string)
	// Extract a value for each key, add it to respData
	for _, k := range keys {
		val, err := rc.db.Get(k).Result()
		if err == nil {
			respData[k] = val
		}
	}
	return respData

}

// Set - Perform set on redis database given a list of pairs
func (rc * RedisDB) Set(items []utils.Pair) bool {
	// TODO: Return Items that were successfully set
	for _, kv := range items {
		_, err := rc.db.Set(kv.Fst.(string), kv.Snd,0).Result()
		if err != nil {
			return false
		}
	}
	return true
}


func (rc *RedisDB) Delete(keys []string) bool {
	// TODO: Return list of successfully deleted items
	_, err := rc.db.Del(keys...).Result()
	if err == nil {
		return false
	}
	return true

}


