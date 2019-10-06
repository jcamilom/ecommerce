package models

import (
	"errors"
	"regexp"
	"strings"

	"github.com/jcamilom/ecommerce/db"
	"github.com/jcamilom/ecommerce/session"
	"golang.org/x/crypto/bcrypt"
)

var (
	// The DB table name for users
	dbUsersTableName = "Users"

	// The DB primary key for users
	dbUsersKeyName = "email"

	// ErrNotFound is returned when a resource cannot be found
	// in the database.
	ErrNotFound = errors.New("models: resource not found")

	// ErrPasswordIncorrect is returned when an invalid password
	// is used when attempting to authenticate a user.
	ErrPasswordIncorrect = errors.New("models: incorrect password provided")

	// ErrEmailRequired is returned when an email address is
	// not provided when creating a user
	ErrEmailRequired = errors.New("models: email address is required")

	// ErrEmailInvalid is returned when an email address provided
	// does not match any of our requirements
	ErrEmailInvalid = errors.New("models: email address is not valid")

	// ErrEmailTaken is returned when an update or create is attempted
	// with an email address that is already in use.
	ErrEmailTaken = errors.New("models: email address is already taken")

	// ErrPasswordRequired is returned when a create is attempted
	// without a user password provided.
	ErrPasswordRequired = errors.New("models: password is required")

	// ErrPasswordTooShort is returned when an update or create is
	// attempted with a user password that is less than 8 characters.
	ErrPasswordTooShort = errors.New("models: password must be at least 8 characters long")

	// ErrNameRequired is returned when a name is not provided
	// when creating a user
	ErrNameRequired = errors.New("models: name is required")
)

const userPwPepper = "secret-random-string"
const sessionKey = "my_secret_key"

// Token expire time in minutes
const sessionExpireTime = 5

// User represents the user model stored in the database
// This is used for user accounts, storing both an email
// address and a password so users can log in and gain
// access to their content.
type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordHash string `json:"password_hash"`
	AccessToken  string `json:"access_token"`
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
	Update(user *User, update interface{}, updateExp string) error
}

// UserService is a set of methods used to manipulate and
// work with the user model
type UserService interface {
	// Authenticate will verify the provided email address and
	// password are correct. If they are correct, the user
	// corresponding to that email will be returned. Otherwise
	// You will receive either:
	// ErrNotFound, ErrPasswordIncorrect, or another error if
	// something goes wrong.
	Authenticate(email, password string) (*User, error)
	Register(user *User) error
	UserDB
}

func NewUserService() UserService {
	udb := newUserDB()
	session := session.NewSessionService(sessionExpireTime, sessionKey)
	uv := newUserValidator(udb)
	return &userService{
		session: session,
		UserDB:  uv,
	}
}

var _ UserService = &userService{}

type userService struct {
	UserDB
	session *session.Session
}

// Register is used to register a new user in the db. Additionally
// an access token is created for the user and returned
func (us *userService) Register(user *User) error {
	err := us.UserDB.Create(user)
	if err != nil {
		return err
	}
	return us.updateToken(user)
}

// Authenticate can be used to authenticate a user with the
// provided email address and password.
// If the email address provided is invalid, this will return
//   nil, ErrNotFound
// If the password provided is invalid, this will return
//   nil, ErrPasswordIncorrect
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
			return nil, ErrPasswordIncorrect
		default:
			return nil, err
		}
	}
	err = us.updateToken(foundUser)
	if err != nil {
		return nil, err
	}

	return foundUser, nil
}

func (us *userService) updateToken(user *User) error {
	token, err := us.session.CreateToken(user.Email)
	if err != nil {
		return err
	}
	user.AccessToken = token
	update := struct {
		AccessToken string `json:":t"`
	}{
		AccessToken: token,
	}
	updateExp := "set access_token = :t"
	return us.UserDB.Update(user, update, updateExp)
}

type userValFunc func(*User) error

func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

var _ UserDB = &userValidator{}

func newUserValidator(udb UserDB) *userValidator {
	return &userValidator{
		UserDB:     udb,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

type userValidator struct {
	UserDB
	emailRegex *regexp.Regexp
}

// ByEmail will normalize the email address before calling
// ByEmail on the UserDB field.
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	if err := runUserValFuncs(&user, uv.normalizeEmail); err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

// Create will create the provided user in the database
func (uv *userValidator) Create(user *User) error {
	err := runUserValFuncs(user,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail,
		uv.normalizeName,
		uv.requiredName,
	)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

// Update update the provided user with the update key and updated expression.
func (uv *userValidator) Update(user *User, update interface{}, updateExp string) error {
	return uv.UserDB.Update(user, update, updateExp)
}

// bcryptPassword will hash a user's password with a
// predefined pepper (userPwPepper) and bcrypt if the
// Password field is not the empty string
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}
	return nil
}

func (uv *userValidator) emailIsAvail(user *User) error {
	_, err := uv.ByEmail(user.Email)
	if err == ErrNotFound {
		// Email address is not taken
		return nil
	}
	if err != nil {
		return err
	}

	return ErrEmailTaken
}

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) normalizeName(user *User) error {
	user.Name = strings.TrimSpace(user.Name)
	return nil
}

func (uv *userValidator) requiredName(user *User) error {
	if user.Name == "" {
		return ErrNameRequired
	}
	return nil
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
	key := userTableQueryKey{
		Email: email,
	}
	found, err := udb.db.GetItem(key, dbUsersTableName, user)
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
	return udb.db.PutItem(dbUsersTableName, user)
}

// updateToken will update the user token field with the data
// specified by the provided user
func (udb *userDB) Update(user *User, update interface{}, updateExp string) error {
	key := userTableQueryKey{
		Email: user.Email,
	}
	return udb.db.UpdateItem(dbUsersTableName, key, update, updateExp)
}

type userTableQueryKey struct {
	Email string `json:"email"`
}
