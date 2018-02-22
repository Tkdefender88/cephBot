package main

import (
	"fmt"

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
	bot.Start()
	<-make(chan struct{})
	return
}
