package configs

import (
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env      string `yaml:"env"`
	Token    string `yaml:"token"`
	Messages `yaml:"messages"`
}

type Messages struct {
	Responses `yaml:"response"`
	Errors    `yaml:"error"`
}

type Responses struct {
	Start    string `yaml:"start"`
	InputUrl string `yaml:"inputUrl"`
}

type Errors struct {
	BadUrl         string `yaml:"badUrl"`
	UnknownCommand string `yaml:"unknownCommand"`
}

func MustConfig() (*Config, error) {
	const op = "congifs.config.MustConfig"

	configPath := "configs/local.yaml"
	if configPath == "" {
		return &Config{}, errors.New(fmt.Sprintf("%s: empty path to local config", op))
	}

	if _, err := os.Stat(configPath); err != nil {
		return &Config{}, errors.New(fmt.Sprintf("%s: config file does not exist", op))
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return &Config{}, err
	}
	return &cfg, nil
}
