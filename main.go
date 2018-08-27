package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Tkdefender88/cephBot/bot"
	"github.com/Tkdefender88/cephBot/config"
)

const token string = ""

func main() {
	err := config.ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//Start up the bot
	goBot, err := bot.Start()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//Wait until CTRL-C or other term signal is recieved to end the program.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("\n-- Good Bye! --")

	//Cleanly close down the discord session
	goBot.Close()
}
