package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockLimiter struct {
	isCalled bool
}

func (m *mockLimiter) Take(ctx context.Context, userID int) error {
	m.isCalled = true
	return nil
}

func TestWaitLimit(t *testing.T) {
	mockLimiter := &mockLimiter{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		userID         int
		request        *http.Request
		expectedStatus int
	}{
		{
			name:           "Authorized Request",
			userID:         123,
			request:        httptest.NewRequest("GET", "/test", nil),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unauthorized Request",
			userID:         0,
			request:        httptest.NewRequest("GET", "/test", nil),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup middleware with mock limiter
			mw := NewLimiterMiddleware(mockLimiter)
			handlerWithMiddleware := mw.WaitLimit(handler)

			// Execute the handler
			rr := httptest.NewRecorder()
			req := tt.request.WithContext(context.WithValue(tt.request.Context(), userIDKey("UserID"), tt.userID))
			handlerWithMiddleware.ServeHTTP(rr, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
