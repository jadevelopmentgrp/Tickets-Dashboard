package api

import (
	"github.com/gin-gonic/gin"
)

func PremiumHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"premium": true,
	})
}
