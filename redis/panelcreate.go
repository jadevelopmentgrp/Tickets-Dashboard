package redis

import (
	"encoding/json"

	"github.com/apex/log"
	database "github.com/jadevelopmentgrp/Tickets-Database"
)

func (c *RedisClient) PublishPanelCreate(settings database.Panel) {
	encoded, err := json.Marshal(settings)
	if err != nil {
		log.Error(err.Error())
		return
	}

	c.Publish(DefaultContext(), "tickets:panel:create", string(encoded))
}
