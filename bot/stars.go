package bot

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var uid = regexp.MustCompile("<@!?[0-9]{18}>")

func init() {
	newCommand("star", 0, false, false, giveStar).setHelp(
		"Args: <user tag> Gives the person a gold star :star:.\n " +
			"It will give a gold star to a user on the server which are" +
			" non-redeemable, meaningless internet brownie points purely" +
			" for popularity. \n If no user tag is given then the bot will" +
			" tell you how many stars you have. By default stars are" +
			" available to everyone in the server but using the `setstars`" +
			" command will switch to admin only mode where only server admins" +
			" may award stars.",
	).add()
	newCommand(
		"setstars",
		discordgo.PermissionManageServer|discordgo.PermissionAdministrator,
		false,
		true,
		setstars,
	).setHelp(
		"Toggles wether or not stars can be awarded by server members or" +
			" only, the admin(s).",
	).add()
}

func giveStar(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) > 1 {
		mention := args[1]
		targetID := strings.Trim(mention, "<>!@")
		//check if argument is a valid user mention
		if !uid.Match([]byte(mention)) {
			s.ChannelMessageSend(m.ChannelID,
				"You must provide a valid user tag")
			return
		}

		//get userID from the mention
		guild, err := guildDetails(m.ChannelID, s)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID,
				"Could not retrieve guild details, no star can be awarded")
			return
		}

		//verify the user isn't awarding themselves a star
		if m.Author.ID == targetID {
			s.ChannelMessageSend(m.ChannelID,
				"You cannot award yourself a star!")
			return
		}

		//verify the user is in the guild
		_, err = s.GuildMember(guild.ID, targetID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID,
				"A problem occurred, make sure the user is in the server")
			return
		}

		server := guildMap.Server[guild.ID]
		userperms, err := permissionDetails(m.Author.ID, m.ChannelID, s)
		hasPerms := userperms&discordgo.PermissionAdministrator > 0
		if !hasPerms {
			s.ChannelMessageSend((m.ChannelID),
				"Only the server owners can award stars")
		}

		//award the star
		if server.Users == nil {
			server.Users = make(map[string]*user)
		}
		targetUsr, ok := server.Users[targetID]
		//check if the user exists the guild struct if not, initilize them
		if !ok {
			targetUsr = &user{
				UserID: targetID,
				Stars:  0,
			}
			server.Users[targetID] = targetUsr
		}
		targetUsr.addStar()
		respMessage := fmt.Sprintf(
			"%s has been awarded a star! :star:, they now have %d stars",
			mention,
			targetUsr.Stars,
		)
		s.ChannelMessageSend(m.ChannelID, respMessage)
		saveServers()
	} else {
		userID := m.Author.ID
		guild, err := guildDetails(m.ChannelID, s)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID,
				"Could not retrieve guild details, no star can be awarded")
		}
		server := guildMap.Server[guild.ID]
		if server.Users == nil {
			server.Users = make(map[string]*user)
		}
		author, ok := server.Users[userID]
		if !ok {
			author = &user{
				UserID: userID,
				Stars:  0,
			}
			server.Users[userID] = author
			saveServers()
		}
		respMessage := fmt.Sprintf("You have %d stars! :star:", author.Stars)
		s.ChannelMessageSend(m.ChannelID, respMessage)
	}
}

func setstars(s *discordgo.Session, m *discordgo.MessageCreate, msg []string) {
	guild, err := guildDetails(m.ChannelID, s)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Could not retrieve server details")
		return
	}

	server := guildMap.Server[guild.ID]
	//toggle the value
	server.LimitStars = !server.LimitStars

	s.ChannelMessageSend(m.ChannelID,
		fmt.Sprintf("Members without admin priveleges %s award stars",
			func() string {
				if server.LimitStars {
					return "cannot"
				}
				return "can"
			}()))
	saveServers()
}

func (u *user) addStar() {
	u.Stars++
}
