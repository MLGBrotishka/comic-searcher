package keyword

import (
	"context"
	"my_app/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockKeywordRepo struct {
	Storage map[string]entity.IdMap
}

func (mr *MockKeywordRepo) Replace(ctx context.Context, invertedIndex map[string]entity.IdMap) error {
	for k, v := range invertedIndex {
		mr.Storage[k] = v
	}
	return nil
}

func (mr *MockKeywordRepo) GetKeywordIds(ctx context.Context, keyword string) (entity.IdMap, error) {
	return mr.Storage[keyword], nil
}

func TestUpdate(t *testing.T) {
	mockRepo := &MockKeywordRepo{
		Storage: make(map[string]entity.IdMap),
	}

	// Sample data to update
	idKeywords := map[int]entity.KeywordMap{
		1: {"example": true, "text": true},
		2: {"example": true, "line": true},
	}

	keywordAdapter := New(mockRepo)

	err := keywordAdapter.Update(context.Background(), idKeywords)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expectedResult := map[string]entity.IdMap{
		"example": {1: true, 2: true},
		"text":    {1: true},
		"line":    {2: true},
	}

	for key, value := range expectedResult {
		assert.Contains(t, mockRepo.Storage, key, "Expected '%s' to be in the storage", key)
		assert.Equal(t, value, mockRepo.Storage[key], "Expected value for key '%s' does not match", key)
	}
}

func TestSearch(t *testing.T) {
	mockRepo := &MockKeywordRepo{
		Storage: map[string]entity.IdMap{
			"example": {1: true, 2: true},
			"text":    {1: true},
			"line":    {2: true},
		},
	}
	keywordAdapter := New(mockRepo)

	keywords := entity.KeywordMap{"example": true, "text": true}

	ids, err := keywordAdapter.Search(context.Background(), keywords)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedIDs := []int{1}
	assert.Equal(t, expectedIDs, ids, "Unexpected IDs returned")

}

type MockComicRepo struct {
	Comics []entity.Comic
}

func (mcr *MockComicRepo) GetAll(ctx context.Context) ([]entity.Comic, error) {
	return mcr.Comics, nil
}

func TestNewLoad(t *testing.T) {
	mockRepo := &MockKeywordRepo{
		Storage: make(map[string]entity.IdMap),
	}
	comics := []entity.Comic{
		{ID: 1, Keywords: entity.KeywordMap{"example": true, "text": true}},
		{ID: 2, Keywords: entity.KeywordMap{"example": true, "line": true}},
	}
	mockComicRepo := &MockComicRepo{
		Comics: comics,
	}

	_, err := NewLoad(mockRepo, mockComicRepo)

	assert.NoError(t, err)
	expectedResult := map[string]entity.IdMap{
		"example": {1: true, 2: true},
		"text":    {1: true},
		"line":    {2: true},
	}

	for key, value := range expectedResult {
		assert.Contains(t, mockRepo.Storage, key, "Expected '%s' to be in the storage", key)
		assert.Equal(t, value, mockRepo.Storage[key], "Expected value for key '%s' does not match", key)
	}
}
