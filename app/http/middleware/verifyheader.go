package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jadevelopmentgrp/Ticket-Dashboard/utils"
)

func VerifyXTicketsHeader(ctx *gin.Context) {
	if ctx.GetHeader("x-tickets") != "true" {
		ctx.AbortWithStatusJSON(400, utils.ErrorStr("Missing x-tickets header"))
	}
}
