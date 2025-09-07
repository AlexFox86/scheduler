package api

import (
	"net/http"
)

// Auth verifies the JWT token
func (h *Handler) Auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := h.auth.Password()
		if len(pass) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Cookie error", http.StatusUnauthorized)
				return
			}

			err = h.auth.ValidateToken(cookie.Value)
			if err != nil {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}

		next(w, r)
	})
}
