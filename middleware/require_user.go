package middleware

import (
	"log"
	"net/http"

	"github.com/jcamilom/ecommerce/models"
)

type RequireUser struct {
	models.UserService
}

func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Println(token)
		// Find user byToken. The call
		next(w, r)
	})
}
