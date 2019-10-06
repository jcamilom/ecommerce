package session

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	// ErrTokenExpired is returned when the session token has expired.
	ErrTokenExpired = errors.New("session: token expired")

	// ErrTokenInvalid is returned when the session token is invalid.
	ErrTokenInvalid = errors.New("session: token invalid")
)

func NewSessionService(expireTime int, tokenKey string) *Session {
	return &Session{
		expireTime: expireTime,
		tokenKey:   []byte(tokenKey),
	}
}

type Session struct {
	expireTime int
	tokenKey   []byte
}

// CreateToken creates a session token for the provided user
func (session *Session) CreateToken(username string) (string, error) {
	// Declare the expiration time of the token
	expirationTime := time.Now().Add(time.Duration(session.expireTime) * time.Minute)
	// Create the JWT claims, which includes the username (email) and expiry time
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(session.tokenKey)
	if err != nil {
		// // If there is an error in creating the JWT return an internal server error
		// w.WriteHeader(http.StatusInternalServerError)
		return "", err
	}
	return tokenString, nil
}

func (session *Session) VerifyToken(token string) (string, error) {
	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// This method will return an error if the token is invalid
	// (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return session.tokenKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", ErrTokenInvalid
		}
		return "", err
	}
	if !tkn.Valid {
		return "", ErrTokenExpired
	}
	// Finally, return the username
	return claims.Username, nil
}

// Claims will help to encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
