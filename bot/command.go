package bot

import (
	"strings"

	"github.com/Tkdefender88/cephBot/config"
	"github.com/bwmarrin/discordgo"
)

//Embeds a the help message of the command c calling the function
func (c command) helpMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Color: config.EmbedColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  c.Name,
				Value: c.Help,
			},
		},
	})
}

func (c command) add() command {
	commandMap[strings.ToLower(c.Name)] = c
	return c
}
