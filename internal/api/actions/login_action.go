package actions

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/domains/user"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/services/token"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/utils"
)

// LoginRequest is a struct that represents the login request
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=5"`
	Password string `json:"password" binding:"required,min=8"`
}

// UserResponse is a struct that represents the user response
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	User        UserResponse
}

// LoginAction is a function that handles the login action
func LoginAction(userRepo user.Repository, tokenManager token.Manager, tokenDuration time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestBody LoginRequest
		if err := ctx.ShouldBindJSON(&requestBody); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		u, err := userRepo.FindByUsername(ctx, requestBody.Username)
		if err != nil {
			if errors.Is(err, user.ErrUserNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": "username or password is incorrect"})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !utils.IsHashPasswordValid(u.Password, requestBody.Password) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "username or password is incorrect"})
			return
		}

		token, err := tokenManager.GenerateToken(u.Username, tokenDuration)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, LoginResponse{
			AccessToken: token,
			User: UserResponse{
				Username:  u.Username,
				CreatedAt: u.CreatedAt,
			},
		})
	}
}
