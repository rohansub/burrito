package db

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
