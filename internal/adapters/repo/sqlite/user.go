package sqlite

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"
	"my_app/internal/entity"
	"my_app/pkg/sqlite"
)

type UserRepo struct {
	*sqlite.Sqlite
}

func NewUser(s *sqlite.Sqlite) *UserRepo {
	return &UserRepo{s}
}

func (r *UserRepo) GetById(ctx context.Context, id int) (entity.User, error) {
	var User entity.User
	var roleBytes []byte

	query := `SELECT id, login, password, salt, role FROM Users WHERE id = $1`
	err := r.Db.QueryRowContext(ctx, query, id).Scan(&User.ID, &User.Login, &User.Password, &User.Salt, &roleBytes)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, fmt.Errorf("UserRepo - GetById - r.Db.QueryRowContext id %d: %w", id, entity.ErrNotFound)
		}
		return entity.User{}, fmt.Errorf("UserRepo - GetById - r.Db.QueryRowContext: %w", err)
	}

	err = gob.NewDecoder(bytes.NewBuffer(roleBytes)).Decode(&User.Admin)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - GetById - gob.NewDecoder.Decode: %w", err)
	}
	return User, nil
}

func (r *UserRepo) GetByLogin(ctx context.Context, login string) (entity.User, error) {
	var User entity.User
	var roleBytes []byte

	query := `SELECT id, login, password, salt, role FROM Users WHERE login = $1`
	err := r.Db.QueryRowContext(ctx, query, login).Scan(&User.ID, &User.Login, &User.Password, &User.Salt, &roleBytes)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, fmt.Errorf("UserRepo - GetById - r.Db.QueryRowContext login %s: %w", login, entity.ErrNotFound)
		}
		return entity.User{}, fmt.Errorf("UserRepo - GetById - r.Db.QueryRowContext: %w", err)
	}

	err = gob.NewDecoder(bytes.NewBuffer(roleBytes)).Decode(&User.Admin)
	if err != nil {
		return entity.User{}, fmt.Errorf("UserRepo - GetById - gob.NewDecoder.Decode: %w", err)
	}
	return User, nil
}

func (r *UserRepo) Store(ctx context.Context, User entity.User) (int, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(User.Admin)
	if err != nil {
		return 0, fmt.Errorf("UserRepo - Store - gob.NewEncoder.Encode: %w", err)
	}

	query := `INSERT INTO Users (login, password, salt, role) VALUES ($1, $2, $3, $4)`
	result, err := r.Db.ExecContext(ctx, query, User.Login, User.Password, User.Salt, buf.Bytes())
	if err != nil {
		return 0, fmt.Errorf("UserRepo - Store - r.Db.ExecContext: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("UserRepo - Store - result.LastInsertId: %w", err)
	}
	return int(id), nil
}
