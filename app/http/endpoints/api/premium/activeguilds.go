package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dbclient "github.com/jadevelopmentgrp/Tickets-Dashboard/database"
	"github.com/jadevelopmentgrp/Tickets-Dashboard/utils"
	"github.com/jadevelopmentgrp/Tickets-Dashboard/utils/types"
	"github.com/jadevelopmentgrp/Tickets-Utilities/permission"
)

type setActiveGuildsBody struct {
	SelectedGuilds types.UInt64StringSlice `json:"selected_guilds"`
}

func SetActiveGuilds(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	var body setActiveGuildsBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorJson(err))
		return
	}

	tx, err := dbclient.Client.BeginTx(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorJson(err))
		return
	}

	defer tx.Rollback(ctx)

	// Validate has admin in each server
	for _, guildId := range body.SelectedGuilds {
		permissionLevel, err := utils.GetPermissionLevel(ctx, guildId, userId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorJson(err))
			return
		}

		if permissionLevel < permission.Admin {
			ctx.JSON(http.StatusForbidden, utils.ErrorStr("Missing permissions in guild %d", guildId))
			return
		}
	}

	ctx.Status(http.StatusNoContent)
}
