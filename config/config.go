package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	//Token is the bots api token
	Token string
	//BotPrefix the default prefix the bot uses to be summoned
	BotPrefix string
	//BotName the name of the bot -Ceph-
	BotName string
	//EmbedColor the default color that the bot embeds it's messages with
	EmbedColor int
	//MentionID is the string that is made when someone mentions the bot
	MentionID = "<@398399749192941568> "

	config *configStruct
)

type configStruct struct {
	Token      string `json:"Token"`
	BotPrefix  string `json:"BotPrefix"`
	BotName    string `json:"BotName"`
	EmbedColor int    `json:"Color"`
}

//ReadConfig reads the config json file and finds the bot Token and the prefix character
func ReadConfig() error {
	fmt.Printf("Reading Config File \n")
	file, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	fmt.Println(string(file))
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	Token = config.Token
	BotPrefix = config.BotPrefix
	BotName = config.BotName
	EmbedColor = config.EmbedColor
	fmt.Println(EmbedColor)
	return nil
}
