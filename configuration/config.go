package configuration

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Credentials struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"credentials"`
		Url string `yaml:"url"`
	} `yaml:"server"`
	Elasticsearch struct {
		Url   string `yaml:"url"`
		Index string `yaml:"index"`
	} `yaml:"elasticsearch"`
}

func GetConfig() Config {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatalf("Read of config failed")
	}
	defer f.Close()

	var config Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Decode of config failed")
	}
	log.Println("Found application configuration, server:", config.Server.Url, "and user:", config.Server.Credentials.Username)
	return config
}
