package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	commandMap = make(map[string]command)

	pingpong = command{"Ping", "\"Pong!\"", false, ping}.add()
	pongping = command{"Pong", "\"Ping!\"", false, pong}.add()
	help     = command{"Help", "", false, msgHelp}.add()
)

//ParseCommand takes in a discord session and a discordgo Message and a message string
//and parses the message string for commands and if found runs the propper commands
func parseCommand(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	//white separate the message to pick out the command parts
	msgList := strings.Fields(message)
	com := toLower(func() string {
		if strings.HasPrefix(message, " ") {
			return " " + msgList[0]
		}
		return msgList[0]
	}())

	if com == toLower(commandMap[com].Name) {
		commandMap[com].Exec(s, m, msgList)
		return
	}
}

func toLower(s string) (r string) {
	return strings.ToLower(s)
}

type command struct {
	Name string
	Help string

	AdminOnly bool

	Exec func(*discordgo.Session, *discordgo.MessageCreate, []string)
}

//Embeds a the help message of the command c calling the function
func (c command) helpMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Color: 0,

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
