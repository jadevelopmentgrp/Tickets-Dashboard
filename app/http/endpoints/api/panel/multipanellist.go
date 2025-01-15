package api

import (
	"context"

	"github.com/gin-gonic/gin"
	dbclient "github.com/jadevelopmentgrp/Ticket-Dashboard/database"
	"github.com/jadevelopmentgrp/Ticket-Dashboard/utils"
	database "github.com/jadevelopmentgrp/Tickets-Database"
	"golang.org/x/sync/errgroup"
)

func MultiPanelList(ctx *gin.Context) {
	type multiPanelResponse struct {
		database.MultiPanel
		Panels []int `json:"panels"`
	}

	guildId := ctx.Keys["guildid"].(uint64)

	multiPanels, err := dbclient.Client.MultiPanels.GetByGuild(ctx, guildId)
	if err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	data := make([]multiPanelResponse, len(multiPanels))
	group, _ := errgroup.WithContext(context.Background())
	for i, multiPanel := range multiPanels {
		i := i
		multiPanel := multiPanel

		data[i] = multiPanelResponse{
			MultiPanel: multiPanel,
		}

		// TODO: Use a join
		group.Go(func() error {
			panels, err := dbclient.Client.MultiPanelTargets.GetPanels(ctx, multiPanel.Id)
			if err != nil {
				return err
			}

			panelIds := make([]int, len(panels))
			for i, panel := range panels {
				panelIds[i] = panel.PanelId
			}

			data[i].Panels = panelIds

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		ctx.JSON(500, utils.ErrorJson(err))
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
		"data":    data,
	})
}
