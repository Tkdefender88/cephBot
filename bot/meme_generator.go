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

//TODO: add a .json file for the memes and a json parser to construct the map
var (
	rollsafe = template{
		name:       "rollsafe",
		filePath:   "./memes/roll_safe.jpg",
		nTextBoxes: 2,
		fontSize:   48,
		wonb:       true,
		textFields: []textField{
			textField{
				x:         351,
				y:         50,
				ax:        0.5,
				ay:        0.2,
				width:     600,
				lineSpace: 1.5,
				align:     gg.AlignCenter,
			},
			textField{
				x:         351,
				y:         345,
				ax:        0.5,
				ay:        0.9,
				width:     600,
				lineSpace: 1.5,
				align:     gg.AlignCenter,
			},
		},
	}.addTemplate()

	brainExpand2 = template{
		name:       "brainexpand2",
		filePath:   "./memes/brain_expand_2.jpg",
		nTextBoxes: 2,
		fontSize:   48,
		wonb:       false,
		textFields: []textField{
			textField{
				x:         170,
				y:         100,
				ax:        0.5,
				ay:        0.5,
				width:     300,
				lineSpace: 1.5,
				align:     gg.AlignCenter,
			},
			textField{
				x:         170,
				y:         360,
				ax:        0.5,
				ay:        0.5,
				width:     300,
				lineSpace: 1.5,
				align:     gg.AlignCenter,
			},
		},
	}.addTemplate()

	brainExpand3 = template{
		name:       "brainexpand3",
		filePath:   "./memes/brain_expand_3.jpg",
		nTextBoxes: 3,
		fontSize:   48,
		wonb:       false,
		textFields: []textField{
			textField{
				x:         214,
				y:         150,
				ax:        0.5,
				ay:        0.5,
				width:     300,
				lineSpace: 1.5,
				align:     gg.AlignCenter,
			},
			textField{
				x:         214,
				y:         430,
				ax:        0.5,
				ay:        0.5,
				width:     300,
				lineSpace: 1.5,
				align:     gg.AlignCenter,
			},
			textField{
				x:         214,
				y:         750,
				ax:        0.5,
				ay:        0.5,
				width:     300,
				lineSpace: 1.5,
				align:     gg.AlignCenter,
			},
		},
	}.addTemplate()

	percy = template{
		name:       "percy",
		filePath:   "memes/percy.png",
		nTextBoxes: 1,
		wonb:       true,
		fontSize:   32,
		textFields: []textField{
			textField{
				x:         30,
				y:         120,
				ax:        0,
				ay:        0.5,
				width:     300,
				lineSpace: 1.2,
				align:     gg.AlignLeft,
			},
		},
	}.addTemplate()

	rainbow = template{
		name:       "rainbowface",
		filePath:   "memes/rainbow_face.jpg",
		nTextBoxes: 1,
		fontSize:   48,
		wonb:       false,
		textFields: []textField{
			textField{
				x:         450,
				y:         460,
				ax:        0.5,
				ay:        0.5,
				width:     850,
				lineSpace: 1.5,
				align:     gg.AlignCenter,
			},
		},
	}.addTemplate()

	pigeon = template{
		name:       "pigeon",
		fontSize:   36,
		filePath:   "memes/isthisapigeon.png",
		nTextBoxes: 3,
		textFields: []textField{
			textField{ //is this a _
				x:         380,
				y:         380,
				ax:        0.5,
				ay:        0.5,
				width:     200,
				lineSpace: 1.5,
				align:     gg.AlignLeft,
			},
			textField{ //The butterfly
				x:         400,
				y:         50,
				ax:        0.5,
				ay:        0.5,
				width:     200,
				lineSpace: 1.5,
				align:     gg.AlignCenter,
			},
			textField{ //Face
				x:         150,
				y:         100,
				ax:        0.5,
				ay:        0.5,
				width:     200,
				lineSpace: 1.5,
				align:     gg.AlignCenter,
			},
		},
		wonb: true,
	}.addTemplate()

	changeMind = template{
		name:       "changemymind",
		fontSize:   32,
		filePath:   "memes/change_my_mind.jpg",
		nTextBoxes: 1,
		wonb:       false,
		textFields: []textField{
			textField{
				x:         320,
				y:         250,
				ax:        0.5,
				ay:        0.5,
				width:     270,
				lineSpace: 1.5,
				align:     gg.AlignCenter,
			},
		},
	}.addTemplate()

	guessIm = template{
		name:       "guessIm",
		fontSize:   36,
		filePath:   "memes/I_guess_Im.jpg",
		nTextBoxes: 0,
		wonb:       false,
	}.addTemplate()
)

//parse the command and generate a meme
func genMeme(s *discordgo.Session, m *discordgo.MessageCreate, msgList []string) {
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
	for i, arg := range arguments { //remove the quotes from each block of text
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
	image, err := gg.LoadImage(t.filePath)
	if err != nil {
		fmt.Println("Could not load the image")
		return image, err
	}
	b := image.Bounds()
	imageWidth := float64(b.Max.X)
	imageHeight := float64(b.Max.Y)

	context := gg.NewContextForImage(image)

	if t.wonb {
		context.SetRGB(1, 1, 1)
	} else {
		context.SetRGB(0, 0, 0)
	}
	if err := context.LoadFontFace("memes/impact.ttf", t.fontSize); err != nil {
		fmt.Println("Could not load the font: ", err)
	}
	fmt.Println(imageWidth, imageHeight)

	var min int
	if t.nTextBoxes < len(args) {
		min = t.nTextBoxes
	} else {
		min = len(args)
	}

	for i := 0; i < min; i++ {
		tf := t.textFields[i]
		context.DrawStringWrapped(
			args[i],
			tf.x,
			tf.y,
			tf.ax,
			tf.ay,
			tf.width,
			tf.lineSpace,
			tf.align,
		)
	}

	return context.Image(), nil
}

//send the meme out after it's been created
func sendMeme(s *discordgo.Session, m *discordgo.MessageCreate, img image.Image) {
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

	fileName := memeAuthor + "'s meme.png"
	meme := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Color: config.EmbedColor,
			Image: &discordgo.MessageEmbedImage{
				URL: "attachment://" + fileName,
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Meme created by: " + memeAuthor,
				IconURL: m.Author.AvatarURL("64"),
			},
		},
		Files: []*discordgo.File{
			&discordgo.File{
				Name:   fileName,
				Reader: bytes.NewReader(buf.Bytes()),
			},
		},
	}
	s.ChannelMessageSendComplex(m.ChannelID, meme)
}

//Gets the meme creators nickname from the guild. If there is no nickname,
//their user name is returned
func getAuthorNick(s *discordgo.Session, m *discordgo.MessageCreate) (memeAuthor string, err error) {
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

//Gets the meme template from the template map based on the name provided by the
//user
func getTemplate(name string) (*template, error) {
	name = strings.ToLower(name)
	if _, ok := templateMap[name]; ok {
		if name == strings.ToLower(templateMap[name].name) {
			t := templateMap[name]
			return t, nil
		}
	}
	return nil, errors.New("Meme template not found")
}

//Adds a template to the map
func (t template) addTemplate() template {
	templateMap[strings.ToLower(t.name)] = &t
	return t
}
