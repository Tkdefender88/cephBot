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

func init() {
	newCommand("moji", 0, false, false, bigMoji).setHelp(
		"Args: [emoji]\n\nPosts a large image of an emoji\nEmoji name must be" +
			" in colon format \n\nExample: `>moji :smile:`\n If a `a` tag is" +
			" given after an emoji then if the emoji given is animated then" +
			" it will display animated and big",
	).add()
}

//bigMoji parses a discord command for an emoji string and attempts to find a
//standard or custom match
func bigMoji(s *discordgo.Session, m *discordgo.MessageCreate,
	msgList []string) {

	//If they don't provide an argument then we don't want to work with them
	if len(msgList) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Must provide an emoji ex. `>moji"+
			" :smile:`")
		return
	}
	//Tries to match the command recieved to a custom emoji tag and ID as
	//discord stores them
	match := emojiRegex.FindStringSubmatch(msgList[1])
	if len(match) == 0 {
		sendEmojiFromFile(s, m, msgList[1])
		return
	}

	//This section builds the emoji url where it is stored on discord also takes
	//into account if it's animated
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

//sendEmojiFromFile attemps to match a given string to a twemoji and send it in
//discord it searches for the twemoji from a folder listing all the standard
//unicode twemojis
func sendEmojiFromFile(s *discordgo.Session, m *discordgo.MessageCreate,
	e string) {

	//turn the string into an emoji file name as set by twemoji
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

//emojiFile takes an emoji string and turns it into magic text of the name
//twemoji gives it does it with math that was figured out by much smarter people
//than me.
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
