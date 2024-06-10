package server

import (
	"encoding/json"
	"net/http"

	"my_app/pkg/logger"
)

type authRoutes struct {
	uc AuthUseCase
	l  logger.Interface
}

func newAuthRoutes(router *http.ServeMux, uc AuthUseCase, mid func(http.HandlerFunc) http.HandlerFunc, l logger.Interface) {
	authRoutes := &authRoutes{uc, l}

	router.HandleFunc("POST /register", mid(authRoutes.register))
	router.HandleFunc("POST /login", mid(authRoutes.login))
}

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"access_token"  binding:"required" example:"auto"`
}

// register обрабатывает запрос на регистрацию нового пользователя
func (a authRoutes) register(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.l.Error(err, "Failed to decode request body")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if req.Login == "" || req.Password == "" {
		a.l.Info("Empty login or password")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	token, err := a.uc.SignUp(r.Context(), req.Login, req.Password)
	if err != nil {
		a.l.Error(err, "Failed to register")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authResponse{token})
}

// login обрабатывает запрос на вход в систему
func (a authRoutes) login(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.l.Error(err, "Failed to decode request body")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if req.Login == "" || req.Password == "" {
		a.l.Info("Empty login or password")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	token, err := a.uc.SignIn(r.Context(), req.Login, req.Password)
	if err != nil {
		a.l.Error(err, "Failed to login")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authResponse{token})
}
