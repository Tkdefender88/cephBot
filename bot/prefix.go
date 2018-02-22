package bot

import (
	"github.com/bwmarrin/discordgo"
)

func setPrefix(s *discordgo.Session, m *discordgo.MessageCreate, msgList []string) {
	guild, err := guildDetails(m.ChannelID, s)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "There was a problem with setting the prefix try again later `1`")
		return
	}

	if len(msgList) < 2 {
		s.ChannelMessageSend(m.ChannelID, "You need to provide a prefix")
		return
	}

	if m.Author.ID != guild.OwnerID && m.Author.ID != juice {
		s.ChannelMessageSend(m.ChannelID, "Only the server owner can perform this magic trick!")
		return
	}
}
