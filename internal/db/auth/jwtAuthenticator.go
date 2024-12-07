package auth

import "github.com/golang-jwt/jwt/v5"

type JWTAuthenticator struct {
	secret string
	aud    string
	iss    string
}

func NewJWTAuthenticator(secret, aud, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{secret: secret, aud: aud, iss: iss}
}

func (a *JWTAuthenticator) GenerateToken(claim jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *JWTAuthenticator) ValidateToken(tokenString string) (*jwt.Token, error) {
	return nil, nil
}
