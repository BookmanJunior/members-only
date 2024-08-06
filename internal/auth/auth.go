package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserClaim struct {
	Id    int
	Admin bool
}

const (
	bearerTokenExp = 10 * time.Minute
)

var secretKey = []byte(os.Getenv("SECRET_TOKEN"))

func CreateToken(currentUser UserClaim) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    currentUser.Id,
		"admin": currentUser.Admin,
		"exp":   time.Now().Add(bearerTokenExp).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	var claims jwt.MapClaims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return claims, err
	}

	claims = token.Claims.(jwt.MapClaims)

	return claims, nil
}
