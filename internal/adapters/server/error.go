package server

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Error string `json:"error" example:"message"`
}

func errorResponse(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response{msg})
}
