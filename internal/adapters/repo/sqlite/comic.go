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

type ComicRepo struct {
	*sqlite.Sqlite
}

func NewComic(s *sqlite.Sqlite) *ComicRepo {
	return &ComicRepo{s}
}

func (r *ComicRepo) Store(ctx context.Context, comic entity.Comic) error {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(comic.Keywords)
	if err != nil {
		return fmt.Errorf("ComicRepo - Store - gob.NewEncoder.Encode: %w", err)
	}

	query := `INSERT INTO comics (id, url, keywords) VALUES ($1, $2, $3)`
	_, err = r.Db.ExecContext(ctx, query, comic.ID, comic.URL, buf.Bytes())
	if err != nil {
		return fmt.Errorf("ComicRepo - Store - r.Db.ExecContext: %w", err)
	}
	return nil
}

func (r *ComicRepo) GetById(ctx context.Context, id int) (entity.Comic, error) {
	var comic entity.Comic
	var keywordsBytes []byte

	query := `SELECT id, url, keywords FROM comics WHERE id = $1`
	err := r.Db.QueryRowContext(ctx, query, id).Scan(&comic.ID, &comic.URL, &keywordsBytes)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Comic{}, fmt.Errorf("ComicRepo - GetById - no comic found with id %d: %w", id, err)
		}
		return entity.Comic{}, fmt.Errorf("ComicRepo - GetById - r.Db.QueryRowContext: %w", err)
	}

	err = gob.NewDecoder(bytes.NewBuffer(keywordsBytes)).Decode(&comic.Keywords)
	if err != nil {
		return entity.Comic{}, fmt.Errorf("ComicRepo - GetById - gob.NewDecoder.Decode: %w", err)
	}
	return comic, nil
}

func (r *ComicRepo) GetAll(ctx context.Context) ([]entity.Comic, error) {
	var comics []entity.Comic
	var keywordsBytes []byte

	query := `SELECT id, url, keywords FROM comics`
	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ComicRepo - GetAll - r.Db.QueryContext: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var comic entity.Comic
		err := rows.Scan(&comic.ID, &comic.URL, &keywordsBytes)
		if err != nil {
			return nil, fmt.Errorf("ComicRepo - GetAll - rows.Scan: %w", err)
		}

		err = gob.NewDecoder(bytes.NewBuffer(keywordsBytes)).Decode(&comic.Keywords)
		if err != nil {
			return nil, fmt.Errorf("ComicRepo - GetById - gob.NewDecoder.Decode: %w", err)
		}

		comics = append(comics, comic)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ComicRepo - GetAll - rows.Err: %w", err)
	}

	return comics, nil
}

func (r *ComicRepo) GetAllIds(ctx context.Context) (entity.IdMap, error) {
	ids := make(entity.IdMap)

	query := `SELECT id FROM comics`
	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ComicRepo - GetAllIds - r.Db.QueryContext: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("ComicRepo - GetAllIds - rows.Scan: %w", err)
		}
		ids[id] = true
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ComicRepo - GetAllIds - rows.Err: %w", err)
	}

	return ids, nil
}
