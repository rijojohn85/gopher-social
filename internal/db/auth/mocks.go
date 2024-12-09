package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TestAuthenticator struct {
	secret string
}

var testClaims = jwt.MapClaims{
	"aud": "test-aud",
	"iss": "test-aud",
	"sub": int64(42),
	"exp": time.Now().Add(time.Hour).Unix(),
}

func NewTestAuthenticator(secret string) *TestAuthenticator {
	return &TestAuthenticator{
		secret: secret,
	}
}

func (a *TestAuthenticator) GenerateToken(claim jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, testClaims)
	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func (a *TestAuthenticator) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.secret), nil
	},
		jwt.WithIssuer("test-aud"),
		jwt.WithAudience("test-aud"),
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
