package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"my_app/internal/adapters/server/middleware"
	"my_app/internal/entity"
	"my_app/pkg/logger"
)

type comicRoutes struct {
	uc ComicUseCase
	l  logger.Interface
}

func newComicRoutes(router *http.ServeMux, uc ComicUseCase, mid func(http.HandlerFunc) http.HandlerFunc, l logger.Interface) {
	routes := &comicRoutes{uc, l}
	roleChecker := middleware.NewRoleMiddleware()
	router.HandleFunc("POST /update", mid(roleChecker.CheckRole(routes.update, true)))
	router.HandleFunc("GET /pics", mid(roleChecker.CheckRole(routes.getPictures, false)))
}

// @Summary     Update comics
// @Description Load new comics
// @ID          update
// @Tags  	    comics
// @Success     200
// @Failure     500
// @Router      /update [post]
func (routes *comicRoutes) update(w http.ResponseWriter, r *http.Request) {
	err := routes.uc.Update(r.Context())
	if err != nil {
		routes.l.Error(err, "http - update")
		errorResponse(w, http.StatusInternalServerError, "internal problems")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

type getPicturesResponse struct {
	Urls []string `json:"urls"       binding:"required"  example:"auto"`
}

// @Summary     Get pictures
// @Description Search most relevant comic and return its picture
// @ID          getPictures
// @Tags  	    comics
// @Accept      query
// @Produce     json
// @Success     200 {object} getPicturesResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /getPictures [get]
func (routes *comicRoutes) getPictures(w http.ResponseWriter, r *http.Request) {
	rawQuery := r.URL.Query().Get("search")
	if len(rawQuery) == 0 {
		routes.l.Error(entity.ErrBadRequest, "http - getPictures")
		errorResponse(w, http.StatusBadRequest, "search should not be empty")
		return
	}
	urls, err := routes.uc.GetPictures(r.Context(), rawQuery)
	if err != nil {
		routes.l.Error(err, "http - getPictures")
		switch errors.Unwrap(err) {
		case entity.ErrBadRequest:
			errorResponse(w, http.StatusNotFound, "provide more information")
		case entity.ErrNotFound:
			errorResponse(w, http.StatusNotFound, "not found")
		default:
			errorResponse(w, http.StatusInternalServerError, "internal problems")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(getPicturesResponse{urls})
}
