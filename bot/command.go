package bot

import (
	"github.com/Tkdefender88/officerDva/config"
	"github.com/bwmarrin/discordgo"
)

type command struct {
	Name string
	Help string

	AdminOnly bool

	Exec func(*discordgo.Session, *discordgo.MessageCreate, []string)
}

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
	commandMap[toLower(c.Name)] = c
	return c
}
