package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	udAPIBase  string = "http://api.urbandictionary.com"
	apiVersion string = "v0"
)

//The result of an API call to the urbandictionary API
type result struct {
	LookupList []struct {
		Word       string `json:"word"`
		Definition string `json:"definition"`
		Example    string `json:"example"`
		Author     string `json:"author"`
		Thumbup    int    `json:"thumbs_up"`
		Thumbdown  int    `json:"thumbs_down"`
	} `json:"list"`
}

func init() {
	newCommand("ud", 0, false, false, udLookup).setHelp(
		"Search words on urban dictionary using `>ud [search word]`",
	).add()
}

func udLookup(s *discordgo.Session, m *discordgo.MessageCreate,
	msgList []string) {

	channel, err := channelDetails(m.ChannelID, s)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not get channel details")
		return
	}

	if !channel.NSFW {
		s.ChannelMessageSend(m.ChannelID,
			"Cannot post Urban dictionary to Non-NSFW channels :(")
		return
	}

	//The search in the url needs to have words joined by plus signs
	search := strings.Join(msgList[1:], "+")
	url := fmt.Sprintf("%s/%s/define?term=%s", udAPIBase, apiVersion, search)
	//Get the json data from the web
	lookupInfo := searchUD(url)
	//Parse the json data for definition author rating and example
	res, err := parseJSONData(lookupInfo)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//If no result is found then send an error message to chat and stop.
	if len(res.LookupList) == 0 {
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Color: embedColor,

			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Error",
					Value: "Definition not found",
				},
			},
		})
		return
	}

	if len(res.LookupList[0].Definition) > 1024 {
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Color: embedColor,

			Fields: []*discordgo.MessageEmbedField{
				{
					Name: "Error",
					Value: "Sorry that definition is too large for a discord" +
						" embed. RIP",
				},
			},
		})
	}
	//Send result as an embeded message
	embedUDresult(s, m, res)
}

func embedUDresult(s *discordgo.Session, m *discordgo.MessageCreate,
	data *result) {

	lookup := data.LookupList[0]
	rating := fmt.Sprintf(":+1:`%d` :-1:`%d`", lookup.Thumbup, lookup.Thumbdown)
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Color: embedColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Definition",
				Value: lookup.Definition,
			},
			{
				Name:  "Example",
				Value: lookup.Example,
			},
			{
				Name:   "Rating",
				Value:  rating,
				Inline: true,
			},
			{
				Name:   "Author",
				Value:  lookup.Author,
				Inline: true,
			},
		},
	})
}

//Sends an HTTP get request to the urban dictionary api and returns the json
//data that it receives as a response
func searchUD(url string) (body []byte) {
	data, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	body, readErr := ioutil.ReadAll(data.Body)
	if readErr != nil {
		fmt.Println(readErr.Error())
		return
	}
	return body
}

func parseJSONData(data []byte) (res *result, err error) {
	jsonErr := json.Unmarshal([]byte(data), &res)
	if jsonErr != nil {
		fmt.Println(jsonErr.Error())
		return res, err
	}
	return res, nil
}
