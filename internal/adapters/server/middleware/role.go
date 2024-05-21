package middleware

import (
	"net/http"
)

type RoleMiddleware struct {
}

func NewRoleMiddleware() *RoleMiddleware {
	return &RoleMiddleware{}
}

func (em *RoleMiddleware) CheckRole(f http.HandlerFunc, isAsmin bool) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		role, ok := req.Context().Value("Role").(bool)
		if isAsmin && !role || !ok {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		f(w, req)
	})
}
