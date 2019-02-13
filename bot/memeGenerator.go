package bot

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"regexp"
	"strings"

	"github.com/Tkdefender88/cephBot/config"
	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
)

var (
	templateMap = make(map[string]*template)
	memeRegex   = regexp.MustCompile(`".*?"`)
	fontPath    = "./memes/unicode.impact.ttf"
)

//a blank template for a meme
type template struct {
	Name       string  //identifies the meme acting as an id for now
	FilePath   string  //defines where the meme template is
	URL        string  `json:"URL,omitempty"`
	NTextBoxes int     //how many text boxes the meme can hold
	FontSize   float64 //how big is the text?
	TextFields []textField
	Wonb       bool //white text on black?
}

//a textfield that gets placed on a meme template
type textField struct {
	X         float64
	Y         float64
	AX        float64
	AY        float64
	Width     float64
	LineSpace float64
	Align     gg.Align
}

func init() {
	newCommand("meme", 0, false, false, genMeme).setHelp(
		"Args: `[template name] <text1> <text2> ... <textn>` \n\n select a" +
			" template and add text to it to make your own meme!\n\nIf no" +
			" meme is specified then a list is DMed to you",
	).add()
	if err := Load("./json/memes.json", &templateMap); err != nil {
		fmt.Println(err)
	}
}

//parse the command and generate a meme
func genMeme(s *discordgo.Session, m *discordgo.MessageCreate,
	msgList []string) {

	if len(msgList) < 2 {
		s.ChannelMessageSend(m.ChannelID,
			"I will send you a list of the templates")
		listTemplates(s, m)
		return
	}

	message := strings.Join(msgList[1:], " ")

	//get the meme template based on the user argument
	templateName := msgList[1]
	template, err := getTemplate(templateName)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	//find all the text arguments in the command based on the regular expression
	arguments := memeRegex.FindAllString(message, -1)
	//remove the quotes from each block of text
	for i, arg := range arguments {
		arguments[i] = arg[1 : len(arg)-1]
	}

	meme, err := addText(template, arguments)
	if err != nil {
		fmt.Println("Could not make the meme")
		return
	}

	sendMeme(s, m, meme)
	if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
		log.Println("failed deleting message: ", err)
	}
}

//Adds text to a meme template
func addText(t *template, args []string) (image.Image, error) {
	image, err := gg.LoadImage(t.FilePath)
	if err != nil {
		fmt.Println("Could not load the image")
		return image, err
	}
	b := image.Bounds()
	imageWidth := float64(b.Max.X)
	imageHeight := float64(b.Max.Y)

	context := gg.NewContextForImage(image)

	if t.Wonb {
		context.SetRGB(1, 1, 1)
	} else {
		context.SetRGB(0, 0, 0)
	}
	if err := context.LoadFontFace("memes/impact.ttf", t.FontSize); err != nil {
		fmt.Println("Could not load the font: ", err)
	}
	fmt.Println(imageWidth, imageHeight)

	var min int
	if t.NTextBoxes < len(args) {
		min = t.NTextBoxes
	} else {
		min = len(args)
	}

	for i := 0; i < min; i++ {
		tf := t.TextFields[i]
		context.DrawStringWrapped(
			args[i],
			tf.X,
			tf.Y,
			tf.AX,
			tf.AY,
			tf.Width,
			tf.LineSpace,
			tf.Align,
		)
	}
	return context.Image(), nil
}

//send the meme out after it's been created
func sendMeme(s *discordgo.Session, m *discordgo.MessageCreate,
	img image.Image) {

	memeAuthor, err := getAuthorNick(s, m)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not resolve authors nickname")
		memeAuthor = m.Author.Username
	}

	//encode the image into a png byte buffer
	buf := &bytes.Buffer{}
	if err := png.Encode(buf, img); err != nil {
		fmt.Println("Could not encode image to png: ", err)
	}

	meme := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Color: config.EmbedColor,
			Image: &discordgo.MessageEmbedImage{
				URL: "attachment://meme.png",
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Meme created by: " + memeAuthor,
				IconURL: m.Author.AvatarURL("64"),
			},
		},
		Files: []*discordgo.File{
			&discordgo.File{
				Name:   "meme.png",
				Reader: bytes.NewReader(buf.Bytes()),
			},
		},
	}
	s.ChannelMessageSendComplex(m.ChannelID, meme)
}

//Gets the meme creators nickname from the guild. If there is no nickname,
//their user name is returned
func getAuthorNick(s *discordgo.Session,
	m *discordgo.MessageCreate) (memeAuthor string, err error) {

	guild, err := guildDetails(m.ChannelID, s)
	if err != nil {
		fmt.Println("Could not resolve guild")
		return "", err
	}

	author, err := s.State.Member(guild.ID, m.Author.ID)
	if err != nil {
		fmt.Println("Could not get author")
		return "", err
	}

	if author.Nick == "" {
		memeAuthor = m.Author.Username
	} else {
		memeAuthor = author.Nick
	}

	return memeAuthor, nil
}

func listTemplates(s *discordgo.Session, m *discordgo.MessageCreate) {
	dm, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not open dm with you")
		fmt.Println("Error opening dm channel: ", err)
		return
	}

	for k, v := range templateMap {
		meme, err := gg.LoadImage(v.FilePath)
		if err != nil {
			fmt.Println("Error loading image: ", err)
			continue
		}

		b := &bytes.Buffer{}
		if err := png.Encode(b, meme); err != nil {
			fmt.Println("Error ecoding image: ", err)
			continue
		}

		e := &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title: k,
				Color: config.EmbedColor,
				Image: &discordgo.MessageEmbedImage{
					URL: "attachment://meme.png",
				},
			},
			Files: []*discordgo.File{
				&discordgo.File{
					Name:        "meme.png",
					ContentType: "image/png",
					Reader:      bytes.NewReader(b.Bytes()),
				},
			},
		}

		s.ChannelMessageSendComplex(dm.ID, e)
	}
}

//Gets the meme template from the template map based on the name provided by the
//user
func getTemplate(name string) (*template, error) {
	name = strings.ToLower(name)
	if _, ok := templateMap[name]; ok {
		if name == strings.ToLower(templateMap[name].Name) {
			t := templateMap[name]
			return t, nil
		}
	}
	return nil, errors.New("Meme template not found")
}

//Adds a template to the map
func (t template) addTemplate() template {
	templateMap[strings.ToLower(t.Name)] = &t
	return t
}
