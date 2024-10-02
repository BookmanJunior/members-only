package auth

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type UserClaim struct {
	Id            int
	Admin         bool
	FileSizeLimit int
}

const (
	bearerTokenExp  = 10 * time.Minute
	refreshTokenExp = 24 * 7 * time.Hour
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
		"id":            currentUser.Id,
		"admin":         currentUser.Admin,
		"fileSizeLimit": currentUser.FileSizeLimit,
		"exp":           time.Now().Add(bearerTokenExp).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateRefreshToken(userId int) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userId,
		"exp": time.Now().Add(refreshTokenExp).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return refreshTokenString, nil
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
