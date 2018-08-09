package bot

import "github.com/bwmarrin/discordgo"

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
	ResultType string `json:"result_type"`
}

//a blank template for a meme
type template struct {
	name       string //identifies the meme acting as an id for now
	filePath   string //defines where the meme template is
	nTextBoxes int    //how many text boxes the meme can hold
	textFields []*textField
	wonb       bool
}

//a textfield that gets placed on a meme template
type textField struct {
	x         float64
	y         float64
	ax        float64
	ay        float64
	width     float64
	lineSpace float64
}

//A command that the bot can recognize
type command struct {
	Name      string
	Help      string
	JuiceOnly bool
	AdminOnly bool
	Handler   func(*discordgo.Session, *discordgo.MessageCreate, []string)
}

type guilds struct {
	Count  int `json:"server_count"`
	Server map[string]*guild
}

type guild struct {
	GuildID       string `json:"guildID"`
	CommandPrefix string `json:"prefix"`
	EmbedColor    int    `json:"embed_color"`
	Kicked        bool   `json:"kicked"`
}

type xkcdComic struct {
	Month      string `json:"month"`
	Num        int    `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Image      string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
}
