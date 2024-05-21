package middleware

import "net/http"

type RoleMiddleware struct {
}

func NewRoleMiddleware() *RoleMiddleware {
	return &RoleMiddleware{}
}

func (em *RoleMiddleware) CheckRole(f http.HandlerFunc, isAsmin bool) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		role := req.Context().Value("role").(bool)
		if !isAsmin || role {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		f(w, req)
	})
}
