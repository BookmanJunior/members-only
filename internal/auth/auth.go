package auth

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

type UserClaim struct {
	Id    int
	Admin bool
}

const (
	bearerTokenExp = 10 * time.Minute
)

var secretKey []byte

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("env file not found")
	}
	secretKey = []byte(os.Getenv("TOKEN_SECRET"))
}

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
