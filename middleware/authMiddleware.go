package middleware

import (
	"net/http"
	"strings"

	"hells/utils"

	"github.com/gorilla/context"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		token := bearerToken[1]
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Set user context for further use
		context.Set(r, "user_id", claims.UserID)
		context.Set(r, "role", claims.Role)

		next.ServeHTTP(w, r)
	})
}

func RBACMiddleware(requiredRole string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			userRole := context.Get(r, "role").(string)

			// Role hierarchy: Admin > Editor > Viewer
			roleHierarchy := map[string]int{
				"Viewer": 1,
				"Editor": 2,
				"Admin":  3,
			}

			if roleHierarchy[userRole] < roleHierarchy[requiredRole] {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
