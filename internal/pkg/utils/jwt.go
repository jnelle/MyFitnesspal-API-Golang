package utils

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func GenJWT(
	claims jwt.MapClaims,
	header map[string]interface{},
	secret []byte,
) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token.Header = header
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DecodeJWT(tokenString string) string {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return nil, nil
	})

	claims := token.Claims.(jwt.MapClaims)

	return fmt.Sprint(claims["sub"])
}
