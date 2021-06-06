package configuration

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Credentails struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"credentials"`
		Url string `yaml:"url"`
	} `yaml:"server"`
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
	return config
}
