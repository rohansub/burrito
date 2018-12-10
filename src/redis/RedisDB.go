package redis

import (
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/rcsubra2/burrito/src/db"
	"github.com/rcsubra2/burrito/src/environment"
	"github.com/rcsubra2/burrito/src/utils"
	re "regexp"
	"time"
)


const (
	singleItem = `(\w+|(?:'\w*'))`
	listOfSingles = `^(?:(?:`+ singleItem + `)\s*,\s*)*$`
	pairItem = `\(\s*`+ singleItem  + `\s*,\s*`+ singleItem +  `\)`
	listOfPairs = `^(?:(?:`+ pairItem + `)\s*,\s*)*$`
)

// RedisDBInterface - interface that is implemented by redis.Client,
// and the RedisMockClient
type RedisDBClientInterface interface {
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(keys ...string) *redis.IntCmd
}


type RedisDatabase struct {
	Client RedisDBClientInterface
}

func NewRedisDatabase(
	isMock bool,
	url    string,
	pswd   string,
) *RedisDatabase {
	var cli RedisDBClientInterface
	if isMock {
		cli = NewMockRedisClient(map[string]string{})
	} else {
		cli = redis.NewClient(&redis.Options{
			Addr:     url,
			Password: pswd, // no password set
			DB:       0,  // use default DB
		})
	}
	return &RedisDatabase {
		Client: cli,
	}
}



func (rd * RedisDatabase) IsCorrectSyntax(fname string, args string) bool {
	if fname == "GET" {
		_, err := NewGetFunction(rd.Client, args)
		if err != nil {
			return false
		}
		return true
	} else if fname == "SET" {
		_, err := NewSetFunction(rd.Client, args)
		if err != nil {
			return false
		}
		return true
	} else if fname == "DEL" {
		_, err := NewDeleteFunction(rd.Client, args)
		if err != nil {
			return false
		}
		return true
	}
	return false
}


func (rd * RedisDatabase) Run(
	fname string,
	args string,
	group environment.EnvironmentGroup,
) (map[string]interface{}, error) {
	// TODO: Refactor this

	if fname == "GET" {
		f, err := NewGetFunction(rd.Client, args)
		if err != nil {
			return nil, err
		}
		return f(group)
	} else if fname == "SET" {
		f, err := NewSetFunction(rd.Client, args)
		if err != nil {
			return nil, err
		}
		return f(group)
	} else if fname == "DEL" {
		f, err := NewDeleteFunction(rd.Client, args)
		if err != nil {
			return nil, err
		}
		return f(group)
	}
	return nil, errors.New("Function " + fname + "not recognized")
}



// Create Get Request
func NewGetFunction(
	client RedisDBClientInterface,
	args string,
)(db.DatabaseFunction, error){
	isGet := re.MustCompile(listOfSingles)
	if isGet.MatchString(args) {
		singleRe := re.MustCompile(singleItem)
		matches := singleRe.FindAllString(args, -1)
		return func(group environment.EnvironmentGroup) (map[string]interface{}, error){
			return Get(matches, client, group), nil
		}, nil
	}
	return nil, errors.New("Invalid argument format for GET: " + args)

}


// Get - Perform get on redis database given a list of strings
func Get(keys []string,
	db RedisDBClientInterface,
	group environment.EnvironmentGroup,
) map[string]interface{} {
	respData := make(map[string]interface{})
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


// Create Set Request
func NewSetFunction(
	client RedisDBClientInterface,
	args string,
)(db.DatabaseFunction, error){

	isSet := re.MustCompile(listOfPairs)
	if isSet.MatchString(args) {
		pairRe := re.MustCompile(pairItem)
		matches := pairRe.FindAllString(args, -1)
		pairs := make([]utils.Pair, len(matches))

		singleRe := re.MustCompile(singleItem)
		for i, m := range matches {
			items := singleRe.FindAllString(m, -1)
			pairs[i] = utils.Pair{
				Fst: items[0],
				Snd: items[1],
			}
		}
		return func(group environment.EnvironmentGroup) (map[string]interface{}, error){
			return nil, Set(pairs, client, group)
		}, nil
	}

	return nil, errors.New("Invalid argument format for GET: " + args)

}


// Set - Perform set on redis database given a list of pairs
func Set(items []utils.Pair, db RedisDBClientInterface, group environment.EnvironmentGroup) error {
	for _, kv := range items {
		fst := extract(kv.Fst.(string), group)
		snd := extract(kv.Snd.(string), group)
		_, err := db.Set(fst, snd,0).Result()
		if err != nil {
			return err
		}
	}
	return nil
}


// Create Delete Request
func NewDeleteFunction(
	client RedisDBClientInterface,
	args string,
)(db.DatabaseFunction, error){
	isGet := re.MustCompile(listOfSingles)
	if isGet.MatchString(args) {
		singleRe := re.MustCompile(singleItem)
		matches := singleRe.FindAllString(args, -1)
		return func(group environment.EnvironmentGroup) (map[string]interface{}, error){
			return nil, Delete(matches, client, group)
		}, nil
	}
	return nil, errors.New("Invalid argument format for GET: " + args)

}

func Delete(keys []string, db RedisDBClientInterface, group environment.EnvironmentGroup) error {
	kvals := make([]string, len(keys))
	for i, k := range keys {
		kvals[i] = extract(k, group)
	}

	_, err := db.Del(kvals...).Result()
	if err != nil {
		return err
	}
	return nil
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
