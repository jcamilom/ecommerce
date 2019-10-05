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

// User represents the user model stored in the database
// This is used for user accounts, storing both an email
// address and a password so users can log in and gain
// access to their content.
type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordHash string `json:"passwordhash"`
}

// UserDB is used to interact with the users database.
//
// For pretty much all single user queries:
// If the user is found, we will return a nil error
// If the user is not found, we will return ErrNotFound
// If there is another error, we will return an error with
// more information about what went wrong. This may not be
// an error generated by the models package.
//
// For single user queries, any error but ErrNotFound should
// probably result in a 500 error.
type UserDB interface {
	// Methods for querying for single users
	ByEmail(email string) (*User, error)
	// Methods for altering users
	Create(user *User) error
}

// UserService is a set of methods used to manipulate and
// work with the user model
type UserService interface {
	// Authenticate will verify the provided email address and
	// password are correct. If they are correct, the user
	// corresponding to that email will be returned. Otherwise
	// You will receive either:
	// ErrNotFound, ErrInvalidPassword, or another error if
	// something goes wrong.
	Authenticate(email, password string) (*User, error)
	UserDB
}

func NewUserService() UserService {
	udb := newUserDB()
	return &userService{
		UserDB: &userValidator{
			UserDB: udb,
		},
	}
}

var _ UserService = &userService{}

type userService struct {
	UserDB
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
func (us *userService) Authenticate(email, password string) (*User, error) {
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

var _ UserDB = &userValidator{}

type userValidator struct {
	UserDB
}

var _ UserDB = &userDB{}

func newUserDB() *userDB {
	db := &db.DB{}
	return &userDB{
		db: db,
	}
}

type userDB struct {
	db *db.DB
}

// ByEmail will look up a user with the provided email.
// If the user is found, we will return a nil error
// If the user is not found, we will return ErrNotFound
func (udb *userDB) ByEmail(email string) (*User, error) {
	user := new(User)
	found, err := udb.db.GetItem(dbKeyName, email, dbTableName, user)
	if err != nil {
		return nil, err
	} else if found == false {
		return nil, ErrNotFound
	} else {
		return user, nil
	}
}

// Create will create the provided user in the database
func (udb *userDB) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return udb.db.PutItem(dbTableName, user)
}
