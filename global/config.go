package global

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var AppConfig *Config

type Config struct {
	System  System  `yaml:"system"`
	LogConf LogConf `yaml:"logConf"`
}

type System struct {
	Environment string
}

func init() {
	data, err := os.ReadFile("./config.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalln(err)
	}
	AppConfig = &config
}
