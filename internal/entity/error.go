package entity

import "errors"

// ErrNotFound not found
var ErrNotFound = errors.New("not found")

var ErrExpired = errors.New("expired")

var ErrWrongCredentials = errors.New("wrong credentials")

var ErrInvalidToken = errors.New("invalid token")

var ErrBadRequest = errors.New("bad request")
