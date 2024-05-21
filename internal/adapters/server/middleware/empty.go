package middleware

import "net/http"

type EmptyMiddleware struct {
}

func NewEmptyMiddleware() *EmptyMiddleware {
	return &EmptyMiddleware{}
}

func (em *EmptyMiddleware) DoNothing(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		f(w, req)
	})
}
