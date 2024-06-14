package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	configFile := "./config_test.yml"

	envVars := map[string]string{
		"APP_NAME":                  "TestApp",
		"APP_VERSION":               "v1.0.0",
		"HTTP_PORT":                 "8080",
		"LOG_LEVEL":                 "info",
		"SQLITE_COMIC_DSN":          "test_comic_dsn",
		"SQLITE_KEYWORD_DSN":        "test_keyword_dsn",
		"SQLITE_USER_DSN":           "test_user_dsn",
		"AUTHORIZER_TOKEN_MAX_TIME": "1h",
		"AUTHORIZER_SECRET":         "secret",
		"FETCHER_URL":               "https://example.com",
		"FETCHER_PARALLEL":          "5",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}

	cfg, err := NewConfig(configFile)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "TestApp", cfg.App.Name)
	assert.Equal(t, "v1.0.0", cfg.App.Version)
	assert.Equal(t, "8080", cfg.HTTP.Port)
	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, "test_comic_dsn", cfg.Sqlite.Comic.Dsn)
	assert.Equal(t, "test_keyword_dsn", cfg.Sqlite.Keyword.Dsn)
	assert.Equal(t, "test_user_dsn", cfg.Sqlite.User.Dsn)
	assert.Equal(t, 1*time.Hour, cfg.Authorizer.TokenMaxTime)
	assert.Equal(t, "secret", cfg.Authorizer.Secret)
	assert.Equal(t, "https://example.com", cfg.Fetcher.URL)
	assert.Equal(t, 5, cfg.Fetcher.Parallel)
}
