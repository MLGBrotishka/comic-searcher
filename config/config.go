package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App     `yaml:"app"`
		HTTP    `yaml:"server"`
		Log     `yaml:"logger"`
		Sqlite  `yaml:"sqlite"`
		Fetcher `yaml:"fetcher"`
	}

	App struct {
		Name    string `yaml:"name"    env:"APP_NAME"`
		Version string `yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `yaml:"log_level"   env:"LOG_LEVEL"`
	}

	Sqlite struct {
		PoolMax ComicDb   `yaml:"pool_max" env:"PG_POOL_MAX"`
		URL     KeywordDb `                env:"PG_URL"`
	}

	ComicDb struct {
		Dsn string `yaml:"pool_max" env:"PG_POOL_MAX"`
	}

	KeywordDb struct {
		Dsn string `yaml:"pool_max" env:"PG_POOL_MAX"`
	}

	Fetcher struct {
		ServerExchange string `yaml:"rpc_server_exchange" env:"RMQ_RPC_SERVER"`
		ClientExchange string `yaml:"rpc_client_exchange" env:"RMQ_RPC_CLIENT"`
		URL            string `                           env:"RMQ_URL"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
