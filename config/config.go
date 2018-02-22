package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	Token      string
	BotPrefix  string
	BotName    string
	EmbedColor int

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
