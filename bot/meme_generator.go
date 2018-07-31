package bot

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/Tkdefender88/cephBot/config"
	"github.com/bwmarrin/discordgo"
	"github.com/golang/freetype"
)

var (
	templateMap = make(map[string]template)
	memeRegex   = regexp.MustCompile(`".*?"`)
	fontPath    = "./memes/unicode.impact.ttf"
)

//TODO: add a .json file for the memes and a json parser to construct the map
var (
	rollsafe = template{
		name:       "rollsafe",
		filePath:   "./memes/roll_safe.jpg",
		nTextBoxes: 2,
		wonb:       true,
		textFields: []*textField{
			&textField{
				x:       50,
				y:       20,
				flow:    "down",
				justify: "center",
			},
			&textField{
				x:       50,
				y:       80,
				flow:    "down",
				justify: "center",
			},
		},
	}.addTemplate()

	brainExpand2 = template{
		name:       "brainexpand2",
		filePath:   "./memes/brain_expand_2.jpg",
		nTextBoxes: 2,
		wonb:       false,
		textFields: []*textField{
			&textField{
				x:       50,
				y:       20,
				flow:    "down",
				justify: "left",
			},
			&textField{
				x:       50,
				y:       40,
				flow:    "down",
				justify: "left",
			},
		},
	}.addTemplate()

	brainExpand3 = template{
		name:       "brainexpand3",
		filePath:   "./memes/brain_expand_3.jpg",
		nTextBoxes: 3,
		wonb:       false,
		textFields: []*textField{
			&textField{
				x:       50,
				y:       20,
				flow:    "down",
				justify: "left",
			},
			&textField{
				x:       50,
				y:       40,
				flow:    "down",
				justify: "left",
			},
			&textField{
				x:       50,
				y:       80,
				flow:    "down",
				justify: "left",
			},
		},
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

	//send the meme
	//sendCompleteMeme(s, m, template)

	//find all the text arguments in the command based on the regular expression
	arguments := memeRegex.FindAllString(message, -1)

	//remove the quotes from each block of text
	for i, arg := range arguments {
		arguments[i] = arg[1 : len(arg)-1]
	}

	r, err := os.Open(template.filePath)
	if err != nil {
		fmt.Print("3")
		s.ChannelMessageSend(m.ChannelID, "r")
		return
	}
	context, err := loadContext(template)
	if err != nil {
		fmt.Print("2")
		s.ChannelMessageSend(m.ChannelID, "context")
		return
	}
	dest, _, err := image.Decode(r)
	if err != nil {
		fmt.Print("1")
		s.ChannelMessageSend(m.ChannelID, "dest")
		return
	}
	if img, ok := dest.(*image.RGBA); ok {
		fmt.Println("ok")
		addLabel(img, template.textFields[0].x, template.textFields[0].y, "Hello World!", context)
		sendCompleteMeme(s, m, img)
	}
}

//TODO: rename this after the old memes feature is completely replaced
//send the meme out after it's been created
//currently being used for testing and just sending out templates
func sendCompleteMeme(s *discordgo.Session, m *discordgo.MessageCreate, img *image.RGBA) {
	/**
	f, err := os.Open(t.filePath)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not open the file for this template")
		return
	}
	defer f.Close()
	*/

	outfile, err := os.Create("./memes/out.png")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "1")
		return
	}
	defer outfile.Close()

	b := bufio.NewWriter(outfile)
	err = png.Encode(b, img)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "2")
		return
	}

	err = b.Flush()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "3")
		return
	}

	f, err := os.Open("./memes/out.png")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "4")
		return
	}
	defer f.Close()

	memeAuthor, err := getAuthor(s, m)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not resolve authors nickname")
		memeAuthor = m.Author.Username
	}

	meme := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Color: config.EmbedColor,
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Meme created by: " + memeAuthor,
				IconURL: m.Author.AvatarURL("64"),
			},
		},
		Files: []*discordgo.File{
			&discordgo.File{
				Name:   "./memes/out.png",
				Reader: f,
			},
		},
	}
	s.ChannelMessageSendComplex(m.ChannelID, meme)
}

//Gets the meme creators nickname from the guild. If there is no nick name,
//their user name is returned
func getAuthor(s *discordgo.Session, m *discordgo.MessageCreate) (memeAuthor string, err error) {
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
	if toLower(name) == toLower(templateMap[name].name) {
		t := templateMap[name]
		return &t, nil
	}
	return nil, errors.New("Meme template not found")
}

//Adds a template to the map
func (t template) addTemplate() template {
	templateMap[toLower(t.name)] = t
	return t
}

//Sets up the context and font to prepare for adding text
func loadContext(t *template) (*freetype.Context, error) {
	//Set up the font
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		fmt.Println("Could not read font file")
		return freetype.NewContext(), errors.New("Could not read font file")
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println("Parse font failed")
		return freetype.NewContext(), errors.New("Parse font failed")
	}

	//set up the context
	fg, bg := image.Black, image.Black
	if t.wonb {
		fg, bg = image.White, image.Black
	}
	rgba := image.NewRGBA(image.Rect(0, 0, 640, 480))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetFont(f)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)

	return c, nil
}

func addLabel(img *image.RGBA, x, y int, label string, c *freetype.Context) {
	c.SetDst(img)
	size := 12.0 //font size in pixels
	pt := freetype.Pt(x, y+int(c.PointToFixed(size)>>6))

	if _, err := c.DrawString(label, pt); err != nil {

	}
}
