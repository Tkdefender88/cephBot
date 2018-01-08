package bot

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	//	"github.com/Tkdefender88/cephBot/config"
	"regexp"
)

var (
	emojiRegex = regexp.MustCompile("<(a)?:.*?:(.*?)>")
)

func bigMoji(s *discordgo.Session, m *discordgo.MessageCreate, msgList []string) {
	if len(msgList) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Must provide an emoji ex. `>moji :smile:`")
		return
	}
	match := emojiRegex.FindStringSubmatch(msgList[1])

	if len(match) == 0 {
		sendEmojiFromFile(s, m, msgList[1])
		return
	}

	url := fmt.Sprintf("https://cdn.discordapp.com/emojis/%s.png", match[2])
	file := "emoji.png"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	defer resp.Body.Close()

	s.ChannelFileSend(m.ChannelID, file, resp.Body)

	s.ChannelMessageDelete(m.ChannelID, m.ID)
}

func sendEmojiFromFile(s *discordgo.Session, m *discordgo.MessageCreate, e string) {
	emoji := emojiFile(e)

	if emoji == "" {
		return
	}

	file, err := os.Open(fmt.Sprintf("emoji/%s.png", emoji))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()

	s.ChannelFileSend(m.ChannelID, "emoji.png", file)

	s.ChannelMessageDelete(m.ChannelID, m.ID)
}

func emojiFile(e string) string {
	found := ""
	file := ""

	for _, r := range e {
		if file != "" {
			file = fmt.Sprintf("%s-%x", file, r)
		} else {
			file = fmt.Sprintf("%x", r)
		}

		if _, err := os.Stat(fmt.Sprintf("emoji/%s.png", file)); err == nil {
			found = file
		} else if file != "" {
			return file
		}
	}
	return found
}
