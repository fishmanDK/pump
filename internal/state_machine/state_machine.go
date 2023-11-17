package state_machine

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"strconv"
)

type StateMachine struct {
	client *redis.Client
}

func (s *StateMachine) GetState(chatId int64) string {
	key := strconv.FormatInt(chatId, 10)
	res := s.client.Get(key)
	return res.Val()
}

func (s *StateMachine) SetupState(chatId int64, state string) error {
	key := strconv.FormatInt(chatId, 10)
	res := s.client.Set(key, state, 0)
	return res.Err()
}

func NewStateMachine() (*StateMachine, error) {
	redisOpts, err := mustRedis()
	if err != nil {
		return &StateMachine{}, err
	}

	opts := redis.Options{
		Addr:     redisOpts.Addr,
		Password: redisOpts.Password,
		DB:       redisOpts.DB,
	}

	return &StateMachine{client: redis.NewClient(&opts)}, nil
}

type RedisStorage struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func mustRedis() (*RedisStorage, error) {
	const op = "state_machine.state_machine.MustRedis"

	configPath := "internal/state_machine/redis.yaml"

	if configPath == "" {
		return &RedisStorage{}, errors.New(fmt.Sprintf("%s: empty path to redis", op))
	}

	_, err := os.Stat(configPath)
	if err != nil {
		return &RedisStorage{}, errors.New(fmt.Sprintf("%s: config file does not exist", op))
	}

	var cfg RedisStorage
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return &RedisStorage{}, err
	}

	return &cfg, nil
}
