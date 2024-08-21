package middlewares

import (
	"api/src/controllers"
	"api/src/security"
	"log"
	"net/http"
	"strconv"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

func Authorize(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := strings.Split(r.Header["Authorization"][0], " ")[1]

		token, err := jwt.Parse(authorization, security.GetJWTSecret)
		if err != nil {
			controllers.Respond(http.StatusUnauthorized, nil, w)
			log.Println(err)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			controllers.Respond(http.StatusUnauthorized, nil, w)
			return
		}

		r.Header.Add(security.UserIDHeader, strconv.Itoa(int(claims["userID"].(float64))))
		next(w, r)
	}
}
