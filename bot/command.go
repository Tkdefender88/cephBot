package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	commandMap = make(map[string]command)
)

//functions that commands execute
type execfunc func(*discordgo.Session, *discordgo.MessageCreate, []string)

//A command that the bot can recognize
type command struct {
	Name string
	Help string

	PermLevel int

	JuiceOnly  bool
	NeedsPerms bool

	Exec execfunc
}

//Embeds a the help message of the command c calling the function
func (c command) helpMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Color: embedColor,
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

func (c command) setHelp(h string) command {
	c.Help = h
	return c
}

//ParseCommand takes in a discord session and a discordgo Message and a message
//string and parses the message string for commands and if found executes the
//related functions
func parseCommand(s *discordgo.Session, m *discordgo.MessageCreate,
	msg string) {

	if len(msg) == 0 {
		return
	}
	//white separate the message to pick out the command parts
	msgList := strings.Fields(msg)
	commandName := func() string {
		if strings.HasPrefix(msgList[0], " ") {
			return " " + msgList[0]
		}
		return msgList[0]
	}()

	if command, ok := commandMap[commandName]; ok &&
		commandName == strings.ToLower(command.Name) {
		perms, err := permissionDetails(m.Author.ID, m.ChannelID, s)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Failed to verify permissions.")
			return
		}

		isJuice := m.Author.ID == juice
		hasPerms := command.PermLevel&perms > 0
		if (!command.JuiceOnly && !command.NeedsPerms) || isJuice || hasPerms {
			command.Exec(s, m, msgList)
			return
		}
		s.ChannelMessageSend(m.ChannelID,
			"You do not possess the power to do this")
		return
	}
}

func newCommand(name string, pl int, jo, np bool, h execfunc) command {
	return command{
		Name:       name,
		PermLevel:  pl,
		JuiceOnly:  jo,
		NeedsPerms: np,
		Exec:       h,
	}
}
