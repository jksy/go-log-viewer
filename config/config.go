package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Inputs []string `json:"inputs"`
}

func LoadConfig(filename string) (*Config, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Print("cant load config file:" + err.Error())
		return nil, err
	}
	var c Config
	json.Unmarshal(contents, &c)
	return &c, nil
}
