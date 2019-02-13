package bot

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Tkdefender88/cephBot/config"
	"github.com/bwmarrin/discordgo"
)

//Meta#3569
//502297955278258189
//miscellaneous commands that aren't large enough to deserve their own file
func init() {
	newCommand("ping", 0, false, false, ping).setHelp("\"Pong!\"").add()
	newCommand("pong", 0, false, false, ping).setHelp("\"Ping!\"").add()
	newCommand("help", 0, false, false, msgHelp).add()
	newCommand("git", 0, false, false, gitHubLink).setHelp(
		"Displays the github link where I'm being developed.",
	).add()
	/*newCommand("request", 0, false, false, featureRequest).setHelp(
	"Requests a feature.").add()*/
	newCommand("report", 0, false, false, bugReport).setHelp(
		"Report a bug.").add()
	newCommand("celebrate", 0, false, false, celebration).setHelp(
		"Starts a celebration!").add()
	newCommand("count", 0, true, true, count).add()

	newCommand("MetaBan", 0, true, false, banUsr).add()
	newCommand("snap", 0, true, false, snap).add()
}

func snap(s *discordgo.Session, m *discordgo.MessageCreate, msg []string) {
	rand.Seed(time.Now().UnixNano())
	if len(msg) < 2 {
		return
	}
	chanID := msg[0]
	guild, err := guildDetails(chanID, s)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for _, m := range guild.Members {
		uID := m.User.ID
		if rand.Intn(100) > 50 {
			if err := s.GuildBanCreateWithReason(guild.ID, uID, "Thanos snapped you", 7); err != nil {
				log.Println(err.Error())
			}
		}
	}
}
func banUsr(s *discordgo.Session, m *discordgo.MessageCreate, msg []string) {
	if len(msg) < 2 {
		return
	}
	uID := msg[1]
	guild, err := guildDetails("501263971890888714", s)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if err := s.GuildBanCreateWithReason(guild.ID, uID, "get yeeted", 7); err != nil {
		log.Println(err.Error())
		return
	}
}

func count(s *discordgo.Session, m *discordgo.MessageCreate, message []string) {
	c, err := channelDetails(countChan, s)
	if err != nil {
		log.Println(err.Error())
		return
	}
	msg, err := msgDetails(c.LastMessageID, countChan, s)
	if err != nil {
		return
	}
	i, err := strconv.Atoi(msg.Content)
	if err != nil {
		return
	}
	s.ChannelMessageSend(countChan, strconv.Itoa(i+1))
}

//ping is basically the hello world test of this whole monstrosity... it worked
//and now we have lots more it sees if the message is a ping it pongs and
//vicea versa
func ping(s *discordgo.Session, m *discordgo.MessageCreate, message []string) {
	if message[0] == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

//msgHelp displays the help messages for the commands of the bot
//if there is no command specified as an argument for the help command then
//it lists all the commands the bot currently knows.
func msgHelp(s *discordgo.Session, m *discordgo.MessageCreate,
	msgList []string) {

	if len(msgList) == 2 {
		if val, ok := commandMap[strings.ToLower(msgList[1])]; ok {
			val.helpMessage(s, m)
			return
		}
		s.ChannelMessageSend(m.ChannelID,
			msgList[1]+" is not a command I know, sorry")
	}

	//create an list of the commands and populate it with each command name
	var commands []string
	for _, val := range commandMap {
		if !val.JuiceOnly {
			commands = append(commands, "`"+val.Name+"`")
		}
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Color: config.EmbedColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: config.BotName,
				Value: strings.Join(commands, ", ") + "\n\n use `" +
					config.BotPrefix + "help [command]` for more details",
			},
		},
	})
}

//gitHubLink this is literally as it reads. It posts a link to my github repo
func gitHubLink(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	s.ChannelMessageSend(
		m.ChannelID,
		"Check out what's under the hood here:"+
			" https://github.com/Tkdefender88/cephBot"+
			"\nLeave a star and make my day! :star:")
}

//celebration is a command just for fun that cheers everyone up and gets the
//party started! woo!
func celebration(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	s.ChannelMessageSend(
		m.ChannelID,
		":sparkles: Woot woot! Time to partay! YAY! :confetti_ball: :tada:",
	)
	s.ChannelMessageDelete(m.ChannelID, m.ID)
}

func featureRequest(s *discordgo.Session, m *discordgo.MessageCreate,
	msgList []string) {
	if len(msgList) < 2 {
		s.ChannelMessageSend(m.ChannelID,
			"If there's a feature you'd like to request please leave a"+
				" message of it after the command. Ex:"+
				" `>request More smiley faces`")
		return
	}

	req := strings.Join(msgList[1:], " ")

	dm, err := s.UserChannelCreate(juice)
	if err != nil {
		fmt.Println("Could not open dm")
	}

	s.ChannelMessageSend(dm.ID,
		":pencil: `"+req+"` requested by: "+m.Author.Username+"#"+
			m.Author.Discriminator)
}

func bugReport(s *discordgo.Session, m *discordgo.MessageCreate,
	msgList []string) {
	if len(msgList) < 2 {
		s.ChannelMessageSend(m.ChannelID,
			"If there is a bug you've found report it with a message please"+
				"Ex: `>report too many smiley faces`")
		return
	}

	report := strings.Join(msgList[1:], " ")

	dm, err := s.UserChannelCreate(juice)
	if err != nil {
		fmt.Println("Could not open dm")
	}

	s.ChannelMessageSend(dm.ID, ":x: `"+report+"` reported by: "+
		m.Author.Username+"#"+m.Author.Discriminator)
}
