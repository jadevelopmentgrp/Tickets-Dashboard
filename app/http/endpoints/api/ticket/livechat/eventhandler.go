package livechat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/jadevelopmentgrp/Tickets-Dashboard/config"
	dbclient "github.com/jadevelopmentgrp/Tickets-Dashboard/database"
	"github.com/jadevelopmentgrp/Tickets-Dashboard/internal/api"
	"github.com/jadevelopmentgrp/Tickets-Dashboard/utils"
)

func (c *Client) HandleEvent(event Event) error {
	switch event.Type {
	case EventTypeAuth:
		var data AuthData
		if err := json.Unmarshal(event.Data, &data); err != nil {
			c.Write(NewErrorMessage("Malformed event payload"))
			_ = c.Ws.Close()
			c.Flush()
			return err
		}

		if err := c.handleAuthEvent(data); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) handleAuthEvent(data AuthData) error {
	if c.Authenticated {
		return api.NewErrorWithMessage(http.StatusBadRequest, errors.New("Already authenticated"), "Already authenticated")
	}

	token, err := jwt.Parse(data.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Conf.Server.Secret), nil
	})
	if err != nil {
		return api.NewErrorWithMessage(http.StatusUnauthorized, err, "Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return api.NewErrorWithMessage(http.StatusUnauthorized, err, "Invalid token data")
	}

	userIdStr, ok := claims["userid"].(string)
	if !ok {
		return api.NewErrorWithMessage(http.StatusUnauthorized, err, "Invalid token data")
	}

	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		return api.NewErrorWithMessage(http.StatusUnauthorized, err, "Invalid token data")
	}

	// Get the ticket
	ticket, err := dbclient.Client.Tickets.Get(context.Background(), c.TicketId, c.GuildId)
	if err != nil {
		return api.NewErrorWithMessage(http.StatusInternalServerError, err, "Error retrieving ticket data")
	}

	if ticket.Id == 0 || ticket.GuildId == 0 || ticket.GuildId != c.GuildId {
		return api.NewErrorWithMessage(http.StatusNotFound, err, "Ticket not found")
	}

	// Verify the user has permissions to be here
	hasPermission, requestErr := utils.HasPermissionToViewTicket(context.Background(), c.GuildId, userId, ticket)
	if requestErr != nil {
		return api.NewErrorWithMessage(http.StatusInternalServerError, err, "Error retrieving permission data")
	}

	if !hasPermission {
		return api.NewErrorWithMessage(http.StatusForbidden, err, "You do not have permission to view this ticket")
	}

	c.Authenticated = true

	c.Write(Event{
		Type: EventTypeAuthenticated,
	})

	return nil
}
