package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

//from Strum355

//channelID and discordgo session pointer
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

func channelDetails(cID string,
	s *discordgo.Session) (*discordgo.Channel, error) {

	channel, err := s.State.Channel(cID)
	if err != nil {
		if err == discordgo.ErrStateNotFound {
			channel, err = s.Channel(cID)
			if err != nil {
				return nil, err
			}
		}
	}
	return channel, nil
}

//authorID, channelID, and Discordgo session pointer
func permissionDetails(aID, cID string, s *discordgo.Session) (int, error) {
	perms, err := s.State.UserChannelPermissions(aID, cID)
	if err != nil {
		if err == discordgo.ErrStateNotFound {
			perms, err = s.UserChannelPermissions(aID, cID)
			if err != nil {
				return 0, err
			}
		}
	}
	return perms, nil
}
