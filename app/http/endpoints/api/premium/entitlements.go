package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	dbclient "github.com/jadevelopmentgrp/Tickets-Dashboard/database"
	"github.com/jadevelopmentgrp/Tickets-Dashboard/utils"
	"github.com/jadevelopmentgrp/Tickets-Dashboard/utils/types"
)

func GetEntitlements(ctx *gin.Context) {
	userId := ctx.Keys["userid"].(uint64)

	res := gin.H{}

	// If it's a multi-server subscription, fetch more data
	guildIds := make([]uint64, 0)
	if err := dbclient.Client.WithTx(ctx, func(tx pgx.Tx) error {
		activeEntitlements, err := dbclient.Client.LegacyPremiumEntitlementGuilds.ListForUser(ctx, tx, userId)
		if err != nil {
			return err
		}

		for _, entitlement := range activeEntitlements {
			guildIds = append(guildIds, entitlement.GuildId)
		}

		return nil
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorJson(err))
		return
	}

	res["selected_guilds"] = types.UInt64StringSlice(guildIds)

	ctx.JSON(http.StatusOK, res)
}
