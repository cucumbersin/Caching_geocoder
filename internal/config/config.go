package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env"`
	DB          `yaml:"db"`
	HTTP_server `yaml:"http_server"`
	API         `yaml:"API"`
}

type DB struct {
	Radis string `yaml:"radis"`
}

type HTTP_server struct {
	IP            string        `yaml:"localhost"`
	Port          int           `yaml:"port"`
	Timeout       time.Duration `yaml:"timeout"`
	Iddle_timeout time.Duration `yaml:"iddle_timeout"`
}

type API struct {
	YandexKey        string `yaml:"yandexKey"`
	GeocodeMapsCoKey string `yaml:"GeocodeMapsCoKey"`
	DaDataApi        `yaml:"daDataApi"`
}
type DaDataApi struct {
	ApiKeyValue    string `yaml:"ApiKeyValue"`
	SecretKeyValue string `yaml:"SecretKeyValue"`
}

func MustLoad() Config {
	configPath := "./config/local.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("config file %s not read error: %s", configPath, err)
	}
	return cfg
}
