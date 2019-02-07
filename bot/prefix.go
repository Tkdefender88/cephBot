package bot

import (
	"github.com/Tkdefender88/cephBot/config"
	"github.com/bwmarrin/discordgo"
)

func init() {
	newCommand(
		"prefix",
		discordgo.PermissionAdministrator|discordgo.PermissionManageServer,
		false,
		true,
		setPrefix,
	).setHelp(
		"`Args: [prefix]`\n\nchanges the prefix that summons me to action" +
			"\nRequires Admin privleges",
	).add()
}

func setPrefix(s *discordgo.Session, m *discordgo.MessageCreate,
	msgList []string) {

	guild, err := guildDetails(m.ChannelID, s)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID,
			"There was a problem with setting the prefix try again later `1`")
		return
	}

	if len(msgList) < 2 {
		s.ChannelMessageSend(m.ChannelID, "You need to provide a prefix")
		return
	}

	if m.Author.ID != guild.OwnerID && m.Author.ID != juice {
		s.ChannelMessageSend(m.ChannelID,
			"Only the server owner can perform this magic trick!")
		return
	}

	server := guildMap.Server[guild.ID]
	server.CommandPrefix = msgList[1]

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Color: config.EmbedColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Prefix Changed",
				Value: "Prefix set to " + msgList[1],
			},
		},
	})
	saveServers()
}
