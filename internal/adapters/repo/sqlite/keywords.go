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

type KeywordRepo struct {
	*sqlite.Sqlite
}

func NewKeyword(s *sqlite.Sqlite) *KeywordRepo {
	return &KeywordRepo{s}
}

func (r *KeywordRepo) GetKeywordIds(ctx context.Context, keyword string) (entity.IdMap, error) {
	query := `SELECT ids FROM keywords WHERE keyword = $1`
	var idsBytes []byte
	err := r.Db.QueryRowContext(ctx, query, keyword).Scan(&idsBytes)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("KeywordRepo - GetKeywordIds - not found keyword %s: %w", keyword, err)
		}
		return nil, fmt.Errorf("KeywordRepo - GetKeywordIds - r.Db.QueryRowContext: %w", err)
	}
	ids := make(entity.IdMap)
	err = gob.NewDecoder(bytes.NewBuffer(idsBytes)).Decode(&ids)
	if err != nil {
		return nil, fmt.Errorf("KeywordRepo - GetKeywordIds - gob.NewDecoder.Decode: %w", err)
	}
	return ids, nil
}

func (r *KeywordRepo) Replace(ctx context.Context, reverseIdx map[string]entity.IdMap) error {
	for keyword := range reverseIdx {
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(reverseIdx[keyword])
		if err != nil {
			return fmt.Errorf("KeywordRepo - Store - gob.NewEncoder.Encode: %w", err)
		}

		query := `Replace INTO keywords (keyword, ids) VALUES ($1, $2)`
		_, err = r.Db.ExecContext(ctx, query, keyword, buf.Bytes())
		if err != nil {
			return fmt.Errorf("KeywordRepo - Store - r.Db.ExecContext: %w", err)
		}
	}
	return nil
}
