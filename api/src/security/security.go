package security

import (
	"api/src/config"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const UserIDHeader = "User-Id"

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func Verify(password string, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func CreateToken(userID uint64) (string, error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"exp":        time.Now().Add(time.Hour).Unix(),
		"userID":     userID,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.JWTSecret))
}

func GetJWTSecret(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("sign method not valid! %s", token.Header["alg"])
	}

	return []byte(config.JWTSecret), nil
}
