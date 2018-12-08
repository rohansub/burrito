package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type DbMeta struct {
	name string
	dbType string
}

// Config - data structure to represent all configuration data
type Config struct {
	name string
	Dbs  map[string]DbMeta
}

func NewConfigFromFile(filename string) (*Config, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config Config
	json.Unmarshal(byteValue, &config)

	return &config, nil

}


