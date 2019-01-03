package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Tkdefender88/cephBot/config"
	"github.com/bwmarrin/discordgo"
)

var (
	guildMap = new(guilds)
	//BotID the bot's ID
	BotID string
	goBot *discordgo.Session
	juice = "146276564726841344" //I am the juice
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
	return loadJSON("servers.json", guildMap)
}

func saveServers() error {
	return saveJSON("servers.json", guildMap)
}

func loadJSON(path string, v interface{}) error {
	file, err := os.OpenFile("json/"+path, os.O_RDONLY, 0600)
	if err != nil {
		fmt.Println("Could not open the file: ", path, err)
		return err
	}
	defer file.Close()
	if err := json.NewDecoder(file).Decode(v); err != nil {
		fmt.Println("Could not load the json", path, err)
		return err
	}
	return nil
}

func saveJSON(path string, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Could not marshal guildMap into json")
	}

	if err := ioutil.WriteFile("json/"+path, bytes, 0600); err != nil {
		fmt.Println("Could not write to file: ", err)
		return err
	}

	return nil
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//fmt.Println(message.Content)
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
}

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
		fmt.Println("Joined new server: ", g.Guild.ID)
	}

	guildMap.Count++

	saveServers()
}
