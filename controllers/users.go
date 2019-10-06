package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jcamilom/ecommerce/models"
)

// NewUsers is used to create a new Users controller
func NewUsers(us models.UserService) *Users {
	return &Users{
		us: us,
	}
}

type Users struct {
	us models.UserService
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
	if ur.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := models.User{
		ID:       "1",
		Name:     ur.Name,
		Email:    ur.Email,
		Password: ur.Password,
	}
	err = u.us.Create(&user)
	switch err {
	case models.ErrEmailRequired, models.ErrEmailInvalid, models.ErrEmailTaken, models.ErrPasswordRequired, models.ErrPasswordTooShort:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&response{
			Message: err.Error(),
		})
	case nil:
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&response{
			Message: fmt.Sprintf("User %v created!", user.Name),
		})
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Login is used to verify the provided email address and
// password and in the db.
//
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ur := new(loginUserRequest)
	err := json.NewDecoder(r.Body).Decode(ur)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if ur.Email == "" || ur.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := u.us.Authenticate(ur.Email, ur.Password)
	switch err {
	case models.ErrNotFound, models.ErrPasswordIncorrect:
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&response{
			Message: "Wrong email - password combination",
		})
	case nil:
		json.NewEncoder(w).Encode(&response{
			Message: fmt.Sprintf("User %v authenticated successfully!", user.Name),
		})
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type createUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type response struct {
	Message string `json:"message"`
}
