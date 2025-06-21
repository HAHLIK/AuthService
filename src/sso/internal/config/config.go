package config

import (
	"os"
	"time"

	"github.com/HAHLIK/AuthService/sso/env"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc" env-requierd:"true"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"8040"`
	Timeout time.Duration `yaml:"timeout" env-default:"10h"`
}

func MustLoad() *Config {
	err := godotenv.Load(env.PATH)
	if err != nil {
		panic(".env file not found")
	}

	path := os.Getenv(env.NAME_CONFIG_PATH)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist")
	}

	var cfg Config

	err = cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic("can't read config file")
	}

	return &cfg
}
