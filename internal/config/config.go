package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Messages struct {
	Responses Responses `yaml:"response"`
	Errors    Errors    `yaml:"error"`
}

type Responses struct {
	Start string `yaml:"start"`
	InputUrl string `yaml:"inputUrl"`
}

type Errors struct {
	UnknownCommand string `yaml:"unknownCommand"`
	BadUrl string `yaml:"badUrl"`
}

type Config struct {
	TelegramToken string `token`
	Env           string `yaml:"env"`

	Messages Messages `yaml:"messages"`
}

func MustConfig() *Config {
	config_path := "configs/local.yaml"

	if config_path == "" {
		log.Fatal("config_path is not set")
	}

	if _, err := os.Stat(config_path); err != nil {
		log.Fatalf("config file does not exist: %s", config_path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(config_path, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err.Error())
	}

	return &cfg
}
