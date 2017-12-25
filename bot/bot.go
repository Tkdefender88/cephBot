package bot

import (
	"fmt"
	"strings"

	"github.com/Tkdefender88/officerDva/config"
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

	goBot.AddHandler(pingPong)
	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running")
}

func pingPong(session *discordgo.Session, message *discordgo.MessageCreate) {
	if strings.HasPrefix(message.Content, config.BotPrefix) {
		if message.Author.ID == BotID {
			return
		}

		if message.Content == config.BotPrefix+"ping" {
			session.ChannelMessageSend(message.ChannelID, "Pong!")
		}

		if message.Content == config.BotPrefix+"pong" {
			session.ChannelMessageSend(message.ChannelID, "Ping!")
		}
	}
}
