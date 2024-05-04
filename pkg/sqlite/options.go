package sqlite

import "time"

type Option func(*Sqlite)

func ConnAttempts(attempts int) Option {
	return func(c *Sqlite) {
		c.connAttempts = attempts
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(c *Sqlite) {
		c.connTimeout = timeout
	}
}
