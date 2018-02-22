package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func guildDetails(cID string, s *discordgo.Session) (*discordgo.Guild, error) {
	channelInGuild, err := s.State.Channel(cID)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	guildDetails, err := s.State.Guild(channelInGuild.GuildID)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return guildDetails, nil
}
