package bot

import (
	"fmt"
	"strings"

	"github.com/Tkdefender88/cephBot/config"
	"github.com/bwmarrin/discordgo"
)

var (
	//BotID the bot's ID
	BotID string
	goBot *discordgo.Session
)

//Start starts the bot session
func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	BotID = u.ID

	goBot.AddHandler(messageCreate)
	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running")
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == BotID {
		return
	}

	if strings.HasPrefix(message.Content, config.BotPrefix) {
		parseCommand(session, message, strings.TrimPrefix(message.Content, config.BotPrefix))
	}
}
