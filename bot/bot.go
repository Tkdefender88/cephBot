package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Tkdefender88/cephBot/config"
	"github.com/bwmarrin/discordgo"
)

var (
	guildMap = new(guilds)
	//BotID the bot's ID
	BotID     string
	goBot     *discordgo.Session
	juice     = "146276564726841344" //I am the juice
	greenmen  = "157896922625998848"
	countChan = "247919139467952128"
)

//Start starts the bot session
func Start() (*discordgo.Session, error) {
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	BotID = u.ID
	if err := loadServers(); err != nil {
		fmt.Println("Could not load server information: ", err)
	}
	goBot.AddHandler(messageCreate)
	goBot.AddHandler(guildJoinEvent)
	err = goBot.Open()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println("Bot is running")
	return goBot, nil
}

func loadServers() error {
	guildMap.Server = make(map[string]*guild)
	return Load("json/servers.json", guildMap)
}

func saveServers() error {
	return Save("json/servers.json", guildMap)
}

//Event handler for message recieve events
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Printf("%s#%s@%s: %s\n", m.Author.Username, m.Author.Discriminator,
		m.ChannelID, m.Content)
	if m.Author.ID == BotID {
		return
	}
	guild, err := guildDetails(m.ChannelID, s)
	var prefix string
	if err != nil {
		fmt.Println("Could not get the guild details")
		prefix = config.BotPrefix
	} else {
		prefix = guildMap.Server[guild.ID].CommandPrefix
	}
	if strings.HasPrefix(m.Content, prefix) {
		parseCommand(s, m, strings.TrimPrefix(m.Content, prefix))
	}
	if strings.HasPrefix(m.Content, config.MentionID) {
		parseCommand(s, m, strings.TrimPrefix(m.Content, config.MentionID))
	}
	if m.ChannelID == countChan {
		i, err := strconv.Atoi(m.Content)
		if err != nil {
			return
		}
		s.ChannelMessageSend(m.ChannelID, strconv.Itoa(i+1))
	}
}

//Event handler for guild join events
func guildJoinEvent(s *discordgo.Session, g *discordgo.GuildCreate) {
	if g.Unavailable {
		fmt.Println("tried to join an unavailable guild: ", g.Guild.ID)
		return
	}

	if _, ok := guildMap.Server[g.Guild.ID]; !ok {
		guildMap.Server[g.Guild.ID] = &guild{
			GuildID:       g.Guild.ID,
			CommandPrefix: config.BotPrefix,
			EmbedColor:    config.EmbedColor,
			Kicked:        false,
		}
		guildMap.Count++
		fmt.Println("Joined new server: ", g.Guild.ID)
	}
	saveServers()
}
