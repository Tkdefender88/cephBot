package bot

import (
	"strings"

	"github.com/Tkdefender88/officerDva/config"
	"github.com/bwmarrin/discordgo"
)

func pong(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	s.ChannelMessageSend(m.ChannelID, "Ping!")
}

func ping(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	s.ChannelMessageSend(m.ChannelID, "Pong!")
}

func msgHelp(s *discordgo.Session, m *discordgo.MessageCreate, msgList []string) {
	if len(msgList) == 2 {
		if val, ok := commandMap[toLower(msgList[1])]; ok {
			val.helpMessage(s, m)
			return
		}
	}

	var commands []string
	for _, val := range commandMap {
		commands = append(commands, "`"+val.Name+"`")
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Color: 0,

		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "OfficerDva",
				Value: strings.Join(commands, ", ") + "\n\n use `" + config.BotPrefix + "help [command]` for more details",
			},
		},
	})
}

func gitHubLink(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	s.ChannelMessageSend(
		m.ChannelID,
		"Check out what's under the hood here: https://github.com/Tkdefender88/discordbot "+
			"\nLeave a star and make Juicetin's day! :star:")
}
