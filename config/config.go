package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App     `yaml:"app"`
		HTTP    `yaml:"http"`
		Log     `yaml:"logger"`
		Sqlite  `yaml:"sqlite"`
		Fetcher `yaml:"fetcher"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level" env:"LOG_LEVEL"`
	}

	Sqlite struct {
		Comic   ComicDb   `env-required:"true" yaml:"comic"`
		Keyword KeywordDb `env-required:"true" yaml:"keyword"`
	}

	ComicDb struct {
		Dsn string `env-required:"true" yaml:"dsn" env:"SQLITE_COMIC_DSN"`
	}

	KeywordDb struct {
		Dsn string `env-required:"true" yaml:"dsn" env:"SQLITE_KEYWORD_DSN"`
	}

	Fetcher struct {
		URL      string `env-required:"true" yaml:"url"      env:"FETCHER_URL"`
		Parallel int    `env-required:"true" yaml:"parallel" env:"FETCHER_PARALLEL"`
	}
)

// NewConfig returns app config.
func NewConfig(path string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
