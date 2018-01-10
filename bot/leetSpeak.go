package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	leetSet = map[string]string{
		"a": "4",
		"e": "3",
		"i": "!",
		"l": "1",
		"o": "0",
		"s": "5",
		"t": "7",
	}
)

//leetSpeak takes in a discord message and converts the characters to leet speak and sends the converted message to chat
func leetSpeak(s *discordgo.Session, m *discordgo.MessageCreate, msgList []string) {
	//Leet string to build
	leetMsg := ""
	//turn their message into one complete string
	msg := strings.Join(msgList[1:], " ")

	//go through every character and convert
	for _, char := range msg {
		if i, ok := leetSet[strings.ToLower(string(char))]; ok {
			leetMsg += i
		} else {
			leetMsg += string(char)
		}
	}
	s.ChannelMessageSend(m.ChannelID, leetMsg)
}
