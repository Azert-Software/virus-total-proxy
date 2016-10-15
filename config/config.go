package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	
)

type Config struct {
	APIkey string
}

func ReadConfig(cfg *Config) {
	b, err := ioutil.ReadFile("config/config.json")

	if err != nil {
		log.Fatal("Error when reading from config")
	}
	
	json.Unmarshal(b, cfg)

	if cfg.APIkey == "" {
		log.Fatal("An APIkey must be supplied")
	}
}
