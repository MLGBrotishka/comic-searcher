package server

import (
	"net/http"

	"my_app/internal/usecase"
	"my_app/pkg/logger"
)

func NewRouter(handler *http.ServeMux, l logger.Interface, uc usecase.Comic) {
	newComicRoutes(handler, uc, l)
}
