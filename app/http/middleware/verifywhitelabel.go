package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jadevelopmentgrp/Tickets-Dashboard/config"
)

func VerifyWhitelabel(isApi bool) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userId := ctx.Keys["userid"].(uint64)

		var isForced bool
		for _, id := range config.Conf.ForceWhitelabel {
			if id == userId {
				isForced = true
				break
			}
		}

		if !isForced {
			if isApi {
				ctx.AbortWithStatusJSON(402, gin.H{
					"success": false,
					"error":   "You must have the whitelabel premium tier",
				})
			} else {
				ctx.Redirect(302, fmt.Sprintf("%s/premium", config.Conf.Server.MainSite))
				ctx.Abort()
			}
			return
		}
	}
}
