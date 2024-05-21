//go:build migrate

package app

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	"github.com/golang-migrate/migrate/source/file"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_defaultAttempts = 20
	_defaultTimeout  = time.Second
)

func init() {
	err := os.MkdirAll("./db", 0777)
	if err != nil {
		log.Fatal(err)
	}
	migrateSqlite3("SQLITE_COMIC_DSN", "./migrations/comics")
	migrateSqlite3("SQLITE_KEYWORD_DSN", "./migrations/keywords")
	migrateSqlite3("SQLITE_USER_DSN", "./migrations/users")
}

func migrateSqlite3(env string, path string) {
	databaseURL, ok := os.LookupEnv(env)
	if !ok || len(databaseURL) == 0 {
		log.Fatalf("migrate: environment variable not declared: %s", env)
	}

	var (
		//attempts = _defaultAttempts
		err error
		m   *migrate.Migrate
	)

	db, err := sql.Open("sqlite3", databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
	}
	fSrc, err := (&file.File{}).Open(path)
	if err != nil {
		log.Fatal(err)
	}

	m, err = migrate.NewWithInstance("file", fSrc, "sqlite3", instance)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	defer m.Close()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Printf("Migrate: no change")
		return
	}

	log.Printf("Migrate: %s up success", env)
}
