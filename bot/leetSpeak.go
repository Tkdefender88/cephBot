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

func init() {
	newCommand("leet", 0, false, false, leetSpeak).setHelp(
		"`Args: [msg]`\nexample: `>leet it's lit fam` converts message to" +
			" !7'5 1!7 f4m",
	).add()
}

//leetSpeak takes in a discord message and converts the characters to leet speak
//and sends the converted message to chat
func leetSpeak(s *discordgo.Session, m *discordgo.MessageCreate,
	msgList []string) {

	if len(msgList) < 2 {
		s.ChannelMessageSend(m.ChannelID,
			"Please include a message for me to modify!")
		return
	}

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
