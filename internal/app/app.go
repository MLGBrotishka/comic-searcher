package app

import (
	"fmt"
	"my_app/config"
	"my_app/internal/adapters/authorizer"
	"my_app/internal/adapters/fetcher"
	"my_app/internal/adapters/hasher"
	"my_app/internal/adapters/keyword"
	repo_sqlite "my_app/internal/adapters/repo/sqlite"
	"my_app/internal/adapters/server"
	"my_app/internal/usecase"
	"my_app/pkg/httpserver"
	"my_app/pkg/logger"
	"my_app/pkg/normalizer"
	"my_app/pkg/sqlite"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	comicSqlite, err := sqlite.NewSqlite(cfg.Sqlite.Comic.Dsn)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - sqlite.NewSqlite: %w", err))
	}
	defer comicSqlite.Close()

	keywordSqlite, err := sqlite.NewSqlite(cfg.Sqlite.Keyword.Dsn)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - sqlite.NewSqlite: %w", err))
	}
	defer keywordSqlite.Close()

	userSqlite, err := sqlite.NewSqlite(cfg.Sqlite.Keyword.Dsn)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - sqlite.NewSqlite: %w", err))
	}
	defer userSqlite.Close()

	// Use case
	comicUseCase := usecase.NewComic(
		repo_sqlite.NewComic(comicSqlite),
		fetcher.NewFetcher(cfg.Fetcher.URL, cfg.Fetcher.Parallel, *l),
		normalizer.New(),
		keyword.New(
			repo_sqlite.NewKeyword(keywordSqlite),
		),
	)

	authUseCase := usecase.NewAuth(
		repo_sqlite.NewUser(userSqlite),
		authorizer.NewAuthorizer(cfg.Authorizer.TokenMaxTime, cfg.Authorizer.Secret),
		hasher.NewHasher(),
	)

	router := http.NewServeMux()
	server.NewRouter(router, comicUseCase, authUseCase, l)
	// HTTP Server
	httpServer := httpserver.New(router, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	l.Info("Server started, waiting for requests...")
	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
