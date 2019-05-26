package guild

import (
	"fmt"
	"github.com/TicketsBot/GoPanel/utils/discord"
	"strconv"
)

func GetGuild(id int) discord.Endpoint {
	return discord.Endpoint{
		RequestType:       discord.GET,
		AuthorizationType: discord.BOT,
		Endpoint:          fmt.Sprintf("/guilds/%s", strconv.Itoa(id)),
	}
}
