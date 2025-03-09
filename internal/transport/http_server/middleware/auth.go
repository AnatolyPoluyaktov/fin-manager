package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const tokenContextKey contextKey = "userToken"

func AuthMiddleware(validToken string) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			token := parts[1]
			if token != validToken {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), tokenContextKey, token)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}

func GetTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(tokenContextKey).(string)
	return token, ok
}
