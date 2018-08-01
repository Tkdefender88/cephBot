package bot

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Tkdefender88/cephBot/config"
	"github.com/bwmarrin/discordgo"
)

var (
	serverMap = new(servers)
	//BotID the bot's ID
	BotID string
	goBot *discordgo.Session
	juice = "146276564726841344" //I am the juice
)

//Start starts the bot session
func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	u, err := goBot.User("@me")
	if err != nil {
		fmt.Println(err.Error())
		return
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
		return
	}
	fmt.Println("Bot is running")
}

func loadServers() error {
	serverMap.Server = make(map[string]*guild)
	return loadJSON("servers.json", serverMap)
}

func saveServers() error {
	return saveJSON("servers.json", serverMap)
}

func loadJSON(path string, v interface{}) error {
	file, err := os.OpenFile("json/"+path, os.O_RDONLY, 0600)
	if err != nil {
		fmt.Println("Could not open the file: ", path, err)
		return err
	}
	if err := json.NewDecoder(file).Decode(v); err != nil {
		fmt.Println("Could not load the json", path, err)
		return err
	}
	return nil
}

func saveJSON(path string, data interface{}) error {
	file, err := os.OpenFile("json/"+path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println("Could not open file", path, err)
		return err
	}
	if err := json.NewEncoder(file).Encode(data); err != nil {
		fmt.Println("Could not save the json", path, err)
		return err
	}
	return nil
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == BotID {
		return
	}
	guild, err := guildDetails(message.ChannelID, session)
	if err != nil {
		fmt.Println("Could not get the guild details")
	}
	prefix := serverMap.Server[guild.ID].CommandPrefix
	if strings.HasPrefix(message.Content, prefix) {
		parseCommand(session, message, strings.TrimPrefix(message.Content, prefix))
	}
}

func guildJoinEvent(s *discordgo.Session, g *discordgo.GuildCreate) {
	if g.Unavailable {
		fmt.Println("tried to join an unavailable guild: ", g.Guild.ID)
		return
	}

	if _, ok := serverMap.Server[g.Guild.ID]; !ok {
		serverMap.Server[g.Guild.ID] = &guild{
			GuildID:       g.Guild.ID,
			CommandPrefix: config.BotPrefix,
			EmbedColor:    config.EmbedColor,
			Kicked:        false,
		}
		fmt.Println("Joined new server: ", g.Guild.ID)
	}

	saveServers()
}
