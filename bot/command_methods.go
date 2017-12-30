package bot

import (
	"github.com/bwmarrin/discordgo"
)

func pong(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	s.ChannelMessageSend(m.ChannelID, "Ping!")
}

func ping(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	s.ChannelMessageSend(m.ChannelID, "Pong!")
}
