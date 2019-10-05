package models

import (
	"errors"

	"github.com/jcamilom/ecommerce/db"
)

var (
	// The DB table name for users
	dbTableName = "Users"

	// The DB primary key for users
	dbKeyName = "Email"

	// ErrNotFound is returned when a resource cannot be found
	// in the database.
	ErrNotFound = errors.New("models: resource not found")
)

func NewUserService() *UserService {
	db := &db.DB{}
	return &UserService{
		db: db,
	}
}

type UserService struct {
	db *db.DB
}

// ByEmail will look up a user with the provided email.
// If the user is found, we will return a nil error
// If the user is not found, we will return ErrNotFound
func (us *UserService) ByEmail(email string) (*User, error) {
	user := new(User)
	found, err := us.db.GetItem(dbKeyName, email, dbTableName, user)
	if err != nil {
		return nil, err
	} else if found == false {
		return nil, ErrNotFound
	} else {
		return user, nil
	}
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
