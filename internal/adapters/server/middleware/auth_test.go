package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"my_app/internal/entity"
)

type mockAuthUseCase struct {
}

func (m *mockAuthUseCase) GetUserFromToken(_ context.Context, token string) (entity.User, error) {
	if token == "invalid_token" {
		return entity.User{}, entity.ErrInvalidToken
	}
	return entity.User{ID: 1, Admin: true}, nil
}

func TestAuthMiddleware(t *testing.T) {

	mock := &mockAuthUseCase{}

	tests := []struct {
		name               string
		authHeader         string
		expectedStatusCode int
	}{
		{
			name:               "Valid Token",
			authHeader:         "valid_token",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid Token",
			authHeader:         "invalid_token",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := NewAuthMiddleware(mock)
			handler := mw.Auth(func(w http.ResponseWriter, r *http.Request) {})
			req, _ := http.NewRequest("GET", "/", nil)
			req.Header.Add("Authorization", tt.authHeader)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatusCode)
			}
		})
	}
}
