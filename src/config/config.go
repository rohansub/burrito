package config

import (
	"encoding/json"
	"github.com/rohansub/burrito/src/db"
	"io/ioutil"
	"os"

	redis "github.com/rohansub/burrito/src/redis"
)


var initializers = map[string]db.DatabaseConstructor{
	"Redis": redis.NewRedisDatabase,
}



type ServerMeta struct {
	Url 	string `json:"url"`
	Password string `json:"password"`
}


type DbMeta struct {
	DbType string     `json:"type"`
	IsMock bool       `json:"is_mock"`
	Server ServerMeta `json:"server"`
}

// Config - data structure to represent all configuration data
type Config struct {
	Name string `json:"name"`
	Databases map[string]DbMeta `json:"databases"`
}

func NewConfigFromFile(filename string) (*Config, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config Config
	err = json.Unmarshal(byteValue, &config)
	return &config, err
}

func (c * Config) CreateDatabaseClients() map[string]db.Database{
	dbmap := make(map[string]db.Database)
	for name, meta := range c.Databases {
		if init, ok := initializers[meta.DbType]; ok {
			dbmap[name] = init(meta.IsMock, meta.Server.Url, meta.Server.Password)
		} else {
			dbmap[name] = nil
		}
	}
	return dbmap
}


