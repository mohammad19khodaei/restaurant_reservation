package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CancelAction is a function that handles the cancel action
func CancelAction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	}
}
