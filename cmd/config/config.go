package config

import (
	. "dcard-assignment/cmd/utils"
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
	CheckError(err)

	err = json.Unmarshal(data, &config)
	CheckError(err)

	return &config
}
