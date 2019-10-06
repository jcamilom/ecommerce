package session

import (
	"time"

	"github.com/dgrijalva/jwt-go"
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

// Claims will help to encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
