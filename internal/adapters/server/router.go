package server

import (
	"net/http"

	"my_app/internal/usecase"
	"my_app/pkg/logger"
)

func NewRouter(handler *http.ServeMux, uc usecase.Comic, auth usecase.Auth, l logger.Interface) {
	newComicRoutes(handler, uc, l)
	newAuthRoutes(handler, auth, l)
}
