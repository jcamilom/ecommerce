package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jcamilom/ecommerce/context"
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
// POST /register
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ur := new(createUserRequest)
	err := json.NewDecoder(r.Body).Decode(ur)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	user := models.User{
		ID:       "1",
		Name:     ur.Name,
		Email:    ur.Email,
		Password: ur.Password,
	}
	err = u.us.Register(&user)
	if err != nil {
		switch err {
		case models.ErrEmailRequired, models.ErrEmailInvalid, models.ErrEmailTaken, models.ErrPasswordRequired, models.ErrPasswordTooShort, models.ErrNameRequired:
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(&messageResponse{
				Message: err.Error(),
			})
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&loginResponse{
		messageResponse{Message: fmt.Sprintf("User %v created!", user.Name)},
		user.AccessToken,
	})
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
	if err != nil {
		switch err {
		case models.ErrNotFound, models.ErrPasswordIncorrect:
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(&messageResponse{
				Message: "Wrong email - password combination",
			})
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	json.NewEncoder(w).Encode(&loginResponse{
		messageResponse{Message: fmt.Sprintf("User %v authenticated successfully!", user.Name)},
		user.AccessToken,
	})
}

// GetFavorites is used to fetch the user favorite products
//
// GET /users/favorites
func (u *Users) GetFavorites(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := context.User(r.Context())
	if user == nil {
		log.Println("Error while fetching the user from the context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user.Favorites)
}

// GetBalance returns the users wallet balance
//
// GET /users/balance
func (u *Users) GetBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := context.User(r.Context())
	if user == nil {
		log.Println("Error while fetching the user from the context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	balance, err := u.us.GetBalance(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(&messageResponse{
		Message: fmt.Sprintf("User balance is '%v' lumens", balance),
	})
}

// GetStoreBalance returns the balance of the e-commerce in the
// stellar network
//
// GET /store/balance
func (u *Users) GetStoreBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := &models.User{
		Wallet: models.Wallet{
			Address: models.StoreStellarAddress,
		},
	}
	balance, err := u.us.GetBalance(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(&messageResponse{
		Message: fmt.Sprintf("The store balance is '%v' lumens", balance),
	})
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

type messageResponse struct {
	Message string `json:"message"`
}

type loginResponse struct {
	messageResponse
	Token string `json:"token"`
}
