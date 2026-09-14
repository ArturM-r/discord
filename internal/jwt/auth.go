package jwt

import (
	"context"
	"fmt"
	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type Claims struct {
	jwtlib.RegisteredClaims
}

type AuthMiddleware struct {
	secret string
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{
		secret: secret,
	}
}

func (m *AuthMiddleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/auth/") {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			JsonError(w, "authorization header required", 401)
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			JsonError(w, "invalid authorization header", 401)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		if tokenString == "" {
			JsonError(w, "token required", 401)
			return
		}

		claims := &Claims{}

		token, err := jwtlib.ParseWithClaims(tokenString, claims, func(token *jwtlib.Token) (any, error) {
			if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.secret), nil
		})
		if err != nil || !token.Valid {
			JsonError(w, "invalid or expired token", 401)
			return
		}

		if claims.Subject == "" {
			JsonError(w, "invalid token subject", 401)
			return
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			JsonError(w, "invalid token subject", 401)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
