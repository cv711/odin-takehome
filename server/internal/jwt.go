package internal

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	secretKey           = "supersecretkey"
	tokenExpirationTime = 15 * time.Minute
	TokenIssuer         = "odin-takehome"
)

func GenerateJWTToken(userId string) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(tokenExpirationTime).Unix(),
		Issuer:    TokenIssuer,
		Subject:   userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ValidateJWTToken(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
