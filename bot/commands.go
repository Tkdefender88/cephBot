package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	commandMap = make(map[string]command)

	pingpong = command{"ping", "\"Pong!\"", false, ping}.add()
	pongping = command{"pong", "\"Ping!\"", false, ping}.add()
	help     = command{"help", "", false, msgHelp}.add()

	celebrate = command{"woot", "starts a celebration!",
		false,
		celebration}.add()

	gitLink = command{
		"git",
		"displays the github link where I'm being developed",
		false,
		gitHubLink}.add()

	memeMachine = command{
		"meme",
		"Args [meme name]\nIf no meme given then a list is sent in pm\n\nPosts" +
			"a dank meme to the chat.",
		false,
		memeMsg}.add()

	urbanLookup = command{
		"ud",
		"Search things on urban dictionary using `>ud [search]`",
		false,
		udLookup}.add()
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

//Wrapper function to save typing. :P
func toLower(s string) (r string) {
	return strings.ToLower(s)
}
