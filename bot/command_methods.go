package bot

import (
	"strings"

	"github.com/Tkdefender88/cephBot/config"
	"github.com/bwmarrin/discordgo"
)

//ping is basically the hello world test of this whole monstrosity... it worked and now we have lots more
//it sees if the message is a ping it pongs and vicea versa
func ping(s *discordgo.Session, m *discordgo.MessageCreate, message []string) {
	if message[0] == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

//msgHelp displays the help messages for the commands of the bot
//if there is no command specified as an argument for the help command then
//it lists all the commands the bot currently knows.
func msgHelp(s *discordgo.Session, m *discordgo.MessageCreate, msgList []string) {
	//checks for an arg specifying a command for help
	if len(msgList) == 2 {
		if val, ok := commandMap[toLower(msgList[1])]; ok {
			//once the string is matched to a command display that commands help string
			val.helpMessage(s, m)
			return
		}
		//if the string is not matched then we inform them and then print the list of commands
		s.ChannelMessageSend(m.ChannelID, msgList[1]+" is not a command I know, sorry")
	}
	//create an list of the commands and populate it with each command name
	var commands []string
	for _, val := range commandMap {
		commands = append(commands, "`"+val.Name+"`") //the back tics create code blocks in discord markdown
	}

	//SEND THE LIST OF ALL COMMANDS WOO!
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

//gitHubLink this is literally as it reads. It posts a link to my github repo for the is repo to chat
func gitHubLink(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	s.ChannelMessageSend(
		m.ChannelID,
		"Check out what's under the hood here: https://github.com/Tkdefender88/cephBot"+
			"\nLeave a star and make Juicetin's day! :star:")
}

//celebration is a command just for fun that cheers everyone up and gets the party started! woo!
func celebration(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	s.ChannelMessageSend(
		m.ChannelID,
		":sparkles: Woot woot! Time to partay! YAY! :confetti_ball: :tada:",
	)

	s.ChannelMessageDelete(m.ChannelID, m.ID)
}
