package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	newCommand("ban", 8, true, true, BanUsr).setHelp(
		"Ban a user. Syntax `ban <user mention> [ban message]`").add()
	newCommand("kick", 8, true, true, kickUsr).setHelp(
		"Kick a user. Syntax `kick <user mention>").add()
}

func kickUsr(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 || len(m.Mentions) == 0 {
		return
	}

	target := m.Mentions[0]

	g, err := guildDetails(m.ChannelID, s)
	if err != nil {
		return
	}

	if err := s.GuildMemberDelete(g.ID, target.ID); err != nil {
		fmt.Printf("Failed to kick user: %v\n", err)
		return
	}
}

// BanUsr bans a user
func BanUsr(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 || len(m.Mentions) == 0 {
		return
	}

	target := m.Mentions[0]

	g, err := guildDetails(m.ChannelID, s)
	if err != nil {
		return
	}

	banMsg := "Banned"
	if len(args) > 2 {
		banMsg = strings.Join(args[2:], " ")
	}

	if err := s.GuildBanCreateWithReason(g.ID, target.ID, banMsg, 1); err != nil {
		fmt.Printf("Failed to ban user %v\n", err)
	}
}
