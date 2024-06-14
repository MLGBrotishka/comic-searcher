package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDoNothing(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	em := NewEmptyMiddleware()
	handler := em.DoNothing(mockHandler)

	req, _ := http.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
