package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		isAdmin        bool
		shouldBeAdmin  bool
		expectedStatus int
	}{
		{
			name:           "Admin Request",
			isAdmin:        true,
			shouldBeAdmin:  true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Non-Admin Request",
			isAdmin:        false,
			shouldBeAdmin:  true,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw := NewRoleMiddleware()
			handlerWithMiddleware := mw.CheckRole(handler, tt.shouldBeAdmin)

			// Execute the handler
			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			ctx := req.Context()
			ctx = context.WithValue(ctx, roleKey("Role"), tt.isAdmin)
			req = req.WithContext(ctx)

			handlerWithMiddleware.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
