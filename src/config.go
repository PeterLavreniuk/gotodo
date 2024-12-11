package gotodo

import (
	"encoding/json"
	"os"
)

type Config struct {
	MySql SqlConfig
}

type SqlConfig struct {
	DatabaseName string
	UserName     string
	Password     string
}

func ReadConfig() (*Config, error) {
	file, err := os.Open("config.json")

	if err != nil {
		return nil, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
