package models

import (
	"errors"

	"github.com/jcamilom/ecommerce/db"
	"golang.org/x/crypto/bcrypt"
)

var (
	// The DB table name for users
	dbTableName = "Users"

	// The DB primary key for users
	dbKeyName = "email"

	// ErrNotFound is returned when a resource cannot be found
	// in the database.
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidPassword is returned when an invalid password
	// is used when attempting to authenticate a user.
	ErrInvalidPassword = errors.New("models: incorrect password provided")
)

const userPwPepper = "secret-random-string"

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

// Authenticate can be used to authenticate a user with the
// provided email address and password.
// If the email address provided is invalid, this will return
//   nil, ErrNotFound
// If the password provided is invalid, this will return
//   nil, ErrInvalidPassword
// If the email and password are both valid, this will return
//   user, nil
// Otherwise if another error is encountered this will return
//   nil, error
func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}

	return foundUser, nil
}

// Create will create the provided user in the database
func (us *UserService) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return us.db.PutItem(dbTableName, user)
}

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordHash string `json:"passwordhash"`
}
