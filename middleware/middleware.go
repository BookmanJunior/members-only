package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bookmanjunior/members-only/config"
	"github.com/bookmanjunior/members-only/handlers"
	"github.com/bookmanjunior/members-only/internal/auth"
)

type LoggerWriter struct {
	http.ResponseWriter
	codeStatus int
}

func (w *LoggerWriter) WriteHeader(codeStatus int) {
	w.ResponseWriter.WriteHeader(codeStatus)
	w.codeStatus = codeStatus
}

func Logger(a *config.Application, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		loggerWriter := LoggerWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(&loggerWriter, r)

		a.InfoLog.Printf("%s %s %s %d %s %s", r.RemoteAddr, r.Proto, r.Method, loggerWriter.codeStatus, r.URL.RequestURI(), time.Since(start))

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
			handlers.Unauthorized(w, a, errors.New("missing Authorization token"))
			return
		}

		bearerToken := strings.Split(authHeader, " ")[1]

		claims, err := auth.VerifyToken(bearerToken)

		if err != nil {
			handlers.Unauthorized(w, a, err)
			return
		}

		currentUser := auth.UserClaim{
			Id:            int(claims["id"].(float64)),
			Admin:         claims["admin"].(bool),
			FileSizeLimit: int(claims["fileSizeLimit"].(float64)),
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "current_user", currentUser)
		r = r.WithContext(ctx)

		fmt.Println("Current claim: ", claims)
		next.ServeHTTP(w, r)
	}
}

func EnableCors(app *config.Application, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://top-members-only-frontend.vercel.app")

		next.ServeHTTP(w, r)
	}
}
