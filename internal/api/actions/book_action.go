package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BookAction is a function that handles the book action
func BookAction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	}
}
