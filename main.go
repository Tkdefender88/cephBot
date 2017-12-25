package main

import (
	"fmt"

	"github.com/Tkdefender88/officerDva/bot"
	"github.com/Tkdefender88/officerDva/config"
)

const token string = ""

var BotID string

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()

	<-make(chan struct{})
	return
}
