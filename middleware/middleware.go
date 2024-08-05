package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/handlers"
	"github.com/bookmanjunior/members-only/internal/auth"
)

func Logger(a *config.Application, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.InfoLog.Printf("%s %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func RecoverPanic(app *config.Application, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close	")
				app.ErrorLog.Printf("%s\n", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func IsAuthorized(a *config.Application, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			handlers.Unauthorized(w, a, errors.New("Missing Authorization token"))
			return
		}

		bearerToken := strings.Split(authHeader, " ")[1]

		claims, err := auth.VerifyToken(bearerToken)

		if err != nil {
			handlers.Unauthorized(w, a, err)
			return
		}

		currentUser := auth.UserClaim{
			Id:    int(claims["id"].(float64)),
			Admin: claims["admin"].(bool),
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "current_user", currentUser)
		r = r.WithContext(ctx)

		fmt.Println("Current claim: ", claims)
		next.ServeHTTP(w, r)
	}
}
