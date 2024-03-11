package config

import (
	"encoding/json"
	"os"
)

type dbConfig struct {
	Addr     string
	Port     int
	Username string
	Name     string
	Password string
}

func GetDbConfig() *dbConfig {
	config := dbConfig{}
	file := "./configs/config.json"
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return &config
}
