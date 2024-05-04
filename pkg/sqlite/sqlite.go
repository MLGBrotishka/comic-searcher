package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

type Sqlite struct {
	connAttempts int
	connTimeout  time.Duration
	Db           *sql.DB
}

func NewSqlite(dsn string, opts ...Option) (*Sqlite, error) {
	sqlite := &Sqlite{
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}
	for _, opt := range opts {
		opt(sqlite)
	}
	var err error
	for sqlite.connAttempts > 0 {
		sqlite.Db, err = sql.Open("sqlite3", dsn)
		if err == nil {
			break
		}
		log.Printf("Sqlite is trying to connect, attempts left: %d", sqlite.connAttempts)
		time.Sleep(sqlite.connTimeout)
		sqlite.connAttempts--
	}
	if err != nil {
		return nil, fmt.Errorf("sqlite - NewSqlite - connAttempts == 0: %w", err)
	}
	return sqlite, nil
}

func (s *Sqlite) Close() {
	if s.Db != nil {
		s.Db.Close()
	}
}
