package actions

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/domains/user"
	"github.com/mohammad19khodaei/restaurant_reservation/internal/utils"
)

// RegisterUserRequest is the request body for the register user action
type RegisterUserRequest struct {
	Username string `json:"username" binding:"required,min=5"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// RegisterUserAction is the action for registering a user
func RegisterUserAction(userRepo user.Repository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var requestBody RegisterUserRequest
		if err := ctx.ShouldBindJSON(&requestBody); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPassword, err := utils.HashPassword(requestBody.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash the provided password"})
			return
		}

		u := &user.User{
			Username: requestBody.Username,
			Password: hashedPassword,
		}
		err = userRepo.Register(ctx, u)
		if err != nil {
			if errors.Is(err, user.ErrUsernameAlreadyExists) {
				ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not register the user"})
			return
		}

		ctx.JSON(http.StatusCreated, UserResponse{
			ID:       u.ID,
			Username: u.Username,
		})
	}
}
