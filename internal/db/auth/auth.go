package auth

import "github.com/golang-jwt/jwt/v5"

type Authenticator interface {
	GenerateToken(claim jwt.Claims) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
}
