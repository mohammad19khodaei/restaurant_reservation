package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/services/token"
)

const (
	AuthUserIDKey           = "auth_user_id"
	AuthorizationTypeBearer = "Bearer"
)

// AuthMiddleware is a Gin middleware that is used to authenticate the request
func AuthMiddleware(tokenManager token.Manager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || parts[0] != AuthorizationTypeBearer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Authorization header format",
			})
			return
		}

		accessToken := parts[1]
		payload, err := tokenManager.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}
		ctx.Set(AuthUserIDKey, payload.UserID)
		ctx.Next()
	}
}
