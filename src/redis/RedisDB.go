package db

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/rcsubra2/burrito/src/db"
	"github.com/rcsubra2/burrito/src/environment"
	"github.com/rcsubra2/burrito/src/utils"
	re "regexp"
	"time"
)


// RedisDBInterface - interface that is implemented by redis.Client,
// and the RedisMockClient
type RedisDBClientInterface interface {
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(keys ...string) *redis.IntCmd
}


// Create database function, given the client, function name, and the argument string
func NewRedisDatabaseFunction(
	client RedisDBClientInterface,
	fname string,
	args string,
) (db.DatabaseFunction, error) {

	if fname == "GET" {
		return nil, nil
	} else {
		return nil, errors.New("Error: " + fname + " is not a Redis function")
	}
}


// Create Get Request
func CreateGet(client RedisDBClientInterface, args string) (db.DatabaseFunction, error){
	single := `(\w+|(?:'\w*'))`
	isGet := re.MustCompile(`^(?:(?:`+ single + `)\s*,\s*)*$`)
	if isGet.MatchString(args) {
		singleRe := re.MustCompile(single)
		matches := singleRe.FindAllString(args, -1)

		return func(group environment.EnvironmentGroup) (interface{}, error){
			return Get(matches, client, group), nil
		}, nil
	}
	return nil, errors.New("Invalid argument format for GET: " + args)

}


// Get - Perform get on redis database given a list of strings
func Get(keys []string,
	db RedisDBClientInterface,
	group environment.EnvironmentGroup,
) map[string]string {
	respData := make(map[string]string)
	// extract a value for each key, add it to respData
	for _, k := range keys {
		kval := extract(k, group)
		val, err := db.Get(kval).Result()
		if err == nil {
			respData[kval] = val
		}
	}
	return respData

}

// Set - Perform set on redis database given a list of pairs
func Set(items []utils.Pair, db RedisDBClientInterface) bool {
	for _, kv := range items {
		_, err := db.Set(kv.Fst.(string), kv.Snd,0).Result()
		if err != nil {
			return false
		}
	}
	return true
}


func Delete(keys []string, db RedisDBClientInterface) bool {
	_, err := db.Del(keys...).Result()
	if err != nil {
		return false
	}
	return true

}



func extract(key string, group environment.EnvironmentGroup) string {
	strRE := re.MustCompile(`^'(\w*)'$`)
	if strRE.MatchString(key)  {
		// handle string
		matches := strRE.FindStringSubmatch(key)
		return matches[1]
	} else {
		val, ok := group.GetValue(key).(string)
		if ok {
			return val
		}
		return ""
	}
}
