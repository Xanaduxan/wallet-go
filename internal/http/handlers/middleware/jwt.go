package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Xanaduxan/wallet-go/internal/service/auth"
	"github.com/golang-jwt/jwt/v5"
)

func JWT(next http.Handler) http.Handler {
	prefix := "Bearer "

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if !strings.HasPrefix(authHeader, prefix) {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len(prefix):]

		keyFunc := func(token *jwt.Token) (any, error) {
			return auth.JwtSecret, nil
		}

		token, err := jwt.Parse(tokenString, keyFunc)

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		emailValue, ok := claims["email"]
		if !ok {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		email, ok := emailValue.(string)
		if !ok || email == "" {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "email", email)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
