package bot

import (
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

var (
	emojiRegex = regexp.MustCompile("<(a)?:.*?:(.*?)>")
)

//bigMoji parses a discord command for an emoji string and attempts to find a standard or custom match
func bigMoji(s *discordgo.Session, m *discordgo.MessageCreate, msgList []string) {
	//If they don't provide an argument then we don't want to work with them
	if len(msgList) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Must provide an emoji ex. `>moji :smile:`")
		return
	}
	//Tries to match the command recieved to a custom emoji tag and ID as disocrd stores them
	match := emojiRegex.FindStringSubmatch(msgList[1])
	//If the message didn't match the regex then we know it's not a custom server emoji
	//which means we need to try and match their message to a standard unicode emoji
	if len(match) == 0 {
		sendEmojiFromFile(s, m, msgList[1])
		return
	}

	//This section builds the emoji url where it is stored on discord also takes into account if it's animated
	var url string
	file := "emoji"
	if len(msgList) == 3 && msgList[2] == "a" {
		url = fmt.Sprintf("https://cdn.discordapp.com/emojis/%s.gif", match[2])
		file += ".gif"
	} else {
		url = fmt.Sprintf("https://cdn.discordapp.com/emojis/%s.png", match[2])
		file += ".png"
	}

	//Hopefully this is where we get an acutall image to send.
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	defer resp.Body.Close()
	//Send the big emoji we found and delete the message that called for it.
	s.ChannelFileSend(m.ChannelID, file, resp.Body)
	s.ChannelMessageDelete(m.ChannelID, m.ID)
}

//sendEmojiFromFile attemps to match a given string to a twemoji and send it in discord
//it searches for the twemoji from a folder listing all the standard unicode twemojis
func sendEmojiFromFile(s *discordgo.Session, m *discordgo.MessageCreate, e string) {
	//turn the string into an emoji file name as set by twemoji
	emoji := emojiFile(e)
	if emoji == "" {
		return
	}
	//Oh no, we didn't find anything cuz people are dumb and didn't ask for a good emoji.
	if emoji == "" {
		return
	}
	//Opens the emoji forlder and tries to find the file of the emoji image
	file, err := os.Open(fmt.Sprintf("emoji/%s.png", emoji))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()
	//Send the emoji and delete the evidence
	s.ChannelFileSend(m.ChannelID, "emoji.png", file)
	s.ChannelMessageDelete(m.ChannelID, m.ID)
}

//emojiFile takes an emoji string and turns it into magic text of the name twemoji gives it
//does it with math that was figured out by much smarter people than me. I can't explain it...
//don't touch it.
func emojiFile(e string) string {
	found := ""
	file := ""
	//This works... so if it ain't broke don't fix it.
	//unless you're looking at it now because it is broken... then this is irrelevant.
	//I'd like to think this will keep working for a while. *shrug*
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
