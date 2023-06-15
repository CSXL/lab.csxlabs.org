package auth

import (
	jwt "github.com/golang-jwt/jwt/v4"
)

type Authorizer struct {
	secret []byte
}

func NewAuthorizer(secret string) *Authorizer {
	return &Authorizer{
		secret: []byte(secret),
	}
}

func (a *Authorizer) GenerateToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	return token.SignedString(a.secret)
}

func (a *Authorizer) GetTokenClaims(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return a.secret, nil
	})
	if err != nil {
		return nil, err
	}

	return parsedToken.Claims.(jwt.MapClaims), nil
}

func (a *Authorizer) ValidateToken(token string) (bool, error) {
	_, err := a.GetTokenClaims(token)
	if err != nil {
		return false, err
	}

	return true, nil
}
