package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/devjoemedia/scrumpilot-go-api/utils"
)

// ContextKey avoids collisions in context values
type ContextKey string

const (
	ContextUserID    ContextKey = "user_id"
	ContextUserEmail ContextKey = "user_email"
)

// RequireAuth validates the Bearer access token and injects
// user claims into the request context.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.Error(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			utils.Error(w, http.StatusUnauthorized, "Invalid Authorization header format")
			return
		}

		claims, err := utils.ParseAccessToken(parts[1])
		if err != nil {
			switch err {
			case utils.ErrExpiredToken:
				utils.Error(w, http.StatusUnauthorized, "Access token expired")
			default:
				utils.Error(w, http.StatusUnauthorized, "Invalid access token")
			}
			return
		}

		// Inject into context for downstream handlers
		ctx := context.WithValue(r.Context(), ContextUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextUserEmail, claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID retrieves the authenticated user's ID from context.
// Returns 0 if not set (i.e. route is not protected).
func GetUserID(r *http.Request) uint {
	if id, ok := r.Context().Value(ContextUserID).(uint); ok {
		return id
	}
	return 0
}

func GetUserIDAndEmailFromRequest(r *http.Request) (uint, string, error) {
	id, ok := r.Context().Value(ContextUserID).(uint)
	if !ok {
		return 0, "", fmt.Errorf("user id not found in context")
	}

	email, ok := r.Context().Value(ContextUserEmail).(string)
	if !ok {
		return 0, "", fmt.Errorf("email not found in context")
	}

	return id, email, nil
}
