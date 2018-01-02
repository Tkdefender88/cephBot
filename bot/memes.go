package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	memes *Memes
)

func memeMsg(s *discordgo.Session, m *discordgo.MessageCreate, msgList []string) {
	//Initialize meme list

	//Get the file of all the memes
	memeFile, err := ioutil.ReadFile("config/memes.json")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = json.Unmarshal([]byte(memeFile), &memes)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//Test for arguments.
	if len(msgList) < 2 {
		listMemes(s, m)
		v, _ := s.ChannelMessage(m.ChannelID, m.ID)
		s.ChannelMessageDelete(m.ChannelID, v.ID)
		return
	}

	//find which meme was selected and send it.
	memeChoice, err := selectMeme(msgList[1])
	if err != nil {
		fmt.Println(err.Error())
		s.ChannelMessageSend(m.ChannelID, "404 meme not found\n either spell better or ask juicetin to add it.")
		return
	}

	sendMeme(s, m, memeChoice)
}

//Goes through the array of memes and compares the name value and the argument from the user
//If one matches then that meme is returned. otherwise an error is sent back
func selectMeme(msg string) (Meme, error) {
	for _, element := range memes.Memes {
		if toLower(element.Name) == toLower(msg) {
			return element, nil
		}
	}
	return memes.Memes[0], errors.New("Meme wasn't found")
}

func sendMeme(s *discordgo.Session, m *discordgo.MessageCreate, meme Meme) {
	if m != nil {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Color: 0xeeff00,

		Image: &discordgo.MessageEmbedImage{
			URL:    meme.Link,
			Width:  100,
			Height: 100,
		},
	})
}

func listMemes(s *discordgo.Session, m *discordgo.MessageCreate) {
	var names []string
	for _, val := range memes.Memes {
		names = append(names, "`"+val.Name+"`\n")
	}

	dmChannel, err := s.UserChannelCreate(m.Author.ID)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	s.ChannelMessageSendEmbed(dmChannel.ID, &discordgo.MessageEmbed{
		Color: 0xeeff00,

		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "OfficerDva",
				Value: strings.Join(names, " "),
			},
		},
	})
}

type Memes struct {
	Memes []Meme `json:"memes"`
}

type Meme struct {
	Name string `json:"Name"`
	Link string `json:"Link"`
}
