package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Tkdefender88/cephBot/config"
	"github.com/bwmarrin/discordgo"
)

const (
	xkcdAPIBase    = "http://xkcd.com/"
	xkcdAPIEnd     = "/info.0.json"
	xkcdMostRecent = "http://xkcd.com/info.0.json"
)

func getXkcd(s *discordgo.Session, m *discordgo.MessageCreate, msgList []string) {
	var url string
	if len(msgList) < 2 {
		url = xkcdMostRecent
	} else {
		url = xkcdAPIBase + msgList[1] + xkcdAPIEnd
	}

	resp, err := http.Get(url)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not find that comic")
		fmt.Println("Could not find the comic", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Could not read the response body", err)
		return
	}

	comicData := xkcdComic{}
	if err := json.Unmarshal(body, &comicData); err != nil {
		fmt.Println("Error unmarshalling json data", err)
		return
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Color: config.EmbedColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Title",
				Value:  comicData.SafeTitle,
				Inline: true,
			},
			{
				Name:   "Number",
				Value:  strconv.Itoa(comicData.Num),
				Inline: true,
			},
		},
		Image: &discordgo.MessageEmbedImage{
			URL: comicData.Image,
		},
		Description: comicData.Alt,
	})
	if err != nil {
		fmt.Println(err)
	}
}
