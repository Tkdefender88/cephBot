package bot

import "github.com/bwmarrin/discordgo"

//The result of an API call to the urbandictionary API
type result struct {
	LookupList []lookup `json:"list"`
	ResultType string   `json:"result_type"`
}

//gets a lookup from the urbandictionary API
type lookup struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
	Example    string `json:"example"`
	Author     string `json:"author"`
	Thumbup    int    `json:"thumbs_up"`
	Thumbdown  int    `json:"thumbs_down"`
}

//Memes is used to store all the used memes
type Memes struct {
	Memes []Meme `json:"memes"`
}

//Meme is a type that stores a name of the meme and the link to the meme
type Meme struct {
	Name string `json:"Name"`
	Link string `json:"Link"`
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
	x       int
	y       int
	flow    string
	justify string
}

//A command that the bot can recognize
type command struct {
	Name      string
	Help      string
	AdminOnly bool
	Exec      func(*discordgo.Session, *discordgo.MessageCreate, []string)
}

type guild struct {
	GuildID       string `json:"guildID"`
	CommandPrefix string `json:"prefix"`
}
