package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/jcamilom/ecommerce/models"
)

// NewUsers is used to create a new Users controller
func NewUsers(us *models.UserService) *Users {
	return &Users{
		us: us,
	}
}

type Users struct {
	us *models.UserService
}

// Create is used to process the register data. This is used
// to create a new user account.
//
// POST /users
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ur := new(createUserRequest)
	err := json.NewDecoder(r.Body).Decode(ur)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if ur.Email == "" || ur.Name == "" || ur.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := models.User{
		ID:       "1",
		Name:     ur.Name,
		Email:    ur.Email,
		Password: ur.Password,
	}
	if err := u.us.Create(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ur)
	}
}

type createUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
