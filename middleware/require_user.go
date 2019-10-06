package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/jcamilom/ecommerce/models"
	"github.com/jcamilom/ecommerce/session"
)

type RequireUser struct {
	models.UserService
}

func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		authValue := r.Header.Get("Authorization")
		token := checkAuthorizationHeader(authValue)
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		user, err := mw.UserService.Authorize(token)
		if err != nil {
			switch err {
			case session.ErrTokenExpired, session.ErrTokenInvalid:
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(&messageResponse{
					Message: err.Error(),
				})
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		r.Header.Set("Email", user.Email)
		log.Printf("User %v authorized\n", user.Email)
		next(w, r)
	})
}

func checkAuthorizationHeader(value string) string {
	if !(strings.HasPrefix(value, "bearer ")) {
		return ""
	}
	sts := strings.Split(value, "bearer ")
	if len(sts) != 2 {
		return ""
	}
	return sts[1]
}

type messageResponse struct {
	Message string `json:"message"`
}
