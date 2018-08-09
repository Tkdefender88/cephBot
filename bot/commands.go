package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	commandMap = make(map[string]command)
	pingpong   = command{"ping", "\"Pong!\"", false, false, ping}.add()
	pongping   = command{"pong", "\"Ping!\"", false, false, ping}.add()
	help       = command{"help", "", false, false, msgHelp}.add()
	celebrate  = command{
		"woot",
		"starts a celebration!",
		false,
		false,
		celebration}.add()
	gitLink = command{
		"git",
		"displays the github link where I'm being developed",
		false,
		false,
		gitHubLink}.add()
	urbanLookup = command{
		"ud",
		"Search things on urban dictionary using `>ud [search]`",
		false,
		false,
		udLookup}.add()
	bigEmojis = command{
		"moji",
		"Args: [emoji]\n\nPosts a large image of an emoji\nEmoji name must be in colon format" +
			"\n\nExample: `>moji :smile:`\n If a `a` tag is given after an emoji then if the emoji" +
			" given is animated then it will display animated and big",
		false,
		false,
		bigMoji}.add()
	leet = command{
		"leet",
		"`Args: [msg]`\nexample: `>leet it's lit fam` converts message to !7'5 1!7 f4m",
		false,
		false,
		leetSpeak}.add()
	prefix = command{
		"prefix",
		"`Args: [prefix]`\n\nchanges the prefix that summons the bot to action\nRequires Admin privleges",
		false,
		true,
		setPrefix}.add()
	meme = command{
		"meme",
		"Args: `[template name] <text1> <text2> ... <textn>` \n\n select a template and add text to it to make your own meme!",
		false,
		false,
		genMeme}.add()
	xkcd = command{
		"xkcd",
		"Args: `<comic number>`\n fetches an xkcd comic from the interwebs\n" +
			"Given a comic number it will fetch a specific comic. If left empty the current comic is posted",
		false,
		false,
		getXkcd}.add()
)

//ParseCommand takes in a discord session and a discordgo Message and a message string
//and parses the message string for commands and if found runs the propper commands
func parseCommand(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	//white separate the message to pick out the command parts
	msgList := strings.Fields(message)
	com := strings.ToLower(func() string {
		if strings.HasPrefix(message, " ") {
			return " " + msgList[0]
		}
		return msgList[0]
	}())
	if com == strings.ToLower(commandMap[com].Name) {
		commandMap[com].Handler(s, m, msgList)
		return
	}
}
