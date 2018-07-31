package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Tkdefender88/cephBot/config"
	"github.com/bwmarrin/discordgo"
)

var (
	memelist *Memes
)

//memeMsg takes in a message and sees if it matches a meme in the repository if it does sends it
func memeMsg(s *discordgo.Session, m *discordgo.MessageCreate, msgList []string) {
	//Initialize meme list
	//Get the file of all the memes
	memeFile, err := ioutil.ReadFile("config/memes.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = json.Unmarshal([]byte(memeFile), &memelist)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//Test for arguments.
	if len(msgList) < 2 {
		listMemes(s, m)
		s.ChannelMessageDelete(m.ChannelID, m.ID)
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
	for _, element := range memelist.Memes {
		if toLower(element.Name) == toLower(msg) {
			return element, nil
		}
	}
	return memelist.Memes[0], errors.New("Meme wasn't found")
}

//sendMeme takes a selected meme and sends it to the chat and deletes the evidence
func sendMeme(s *discordgo.Session, m *discordgo.MessageCreate, me Meme) {
	if m != nil {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Color: config.EmbedColor,
		Image: &discordgo.MessageEmbedImage{
			URL:    me.Link,
			Width:  100,
			Height: 100,
		},
	})
}

//listMemes sends a dm to the user requesting with a list of all the memes available
func listMemes(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Populate a list with meme names
	var names []string
	for _, val := range memelist.Memes {
		names = append(names, "`"+val.Name+"`\n")
	}

	//get a dm channel ID
	dmChannel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	s.ChannelMessageSendEmbed(dmChannel.ID, &discordgo.MessageEmbed{
		Color: config.EmbedColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  config.BotName,
				Value: strings.Join(names, " "),
			},
		},
	})
}

// //Memes is used to store all the used memes
// type Memes struct {
// 	Memes []Meme `json:"memes"`
// }

// //Meme is a type that stores a name of the meme and the link to the meme
// type Meme struct {
// 	Name string `json:"Name"`
// 	Link string `json:"Link"`
// }
