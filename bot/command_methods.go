package bot

import (
	"strings"

	"github.com/Tkdefender88/cephBot/config"
	"github.com/bwmarrin/discordgo"
)

func ping(s *discordgo.Session, m *discordgo.MessageCreate, message []string) {
	if message[0] == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
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
		Color: config.EmbedColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  config.BotName,
				Value: strings.Join(commands, ", ") + "\n\n use `" + config.BotPrefix + "help [command]` for more details",
			},
		},
	})
}

func gitHubLink(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	s.ChannelMessageSend(
		m.ChannelID,
		"Check out what's under the hood here: https://github.com/Tkdefender88/cephBot"+
			"\nLeave a star and make Juicetin's day! :star:")
}

func celebration(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	s.ChannelMessageSend(
		m.ChannelID,
		":sparkles: Woot woot! Time to partay! YAY! :confetti_ball: :tada:",
	)
	s.ChannelMessageDelete(m.ChannelID, m.ID)
}
