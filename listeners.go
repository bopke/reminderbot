package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"time"
)

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore own messages
	if m.Author.ID == s.State.User.ID {
		return
	}
	// ignore other bots messages
	if m.Author.Bot {
		return
	}
	if handleAdminCommand(s, m) {
		return
	}
	if handleUserCommand(s, m) {
		return
	}
}

func handleAdminCommand(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	args := strings.Split(m.Content, " ")
	if len(args) <= 1 {
		return false
	}
	if !(args[0] == "<@"+s.State.User.ID+">" || args[0] == "<@!"+s.State.User.ID+">") {
		return false
	}
	switch args[1] {
	case "setPermittedRole":
		break
	default:
		return false
	}
	return true
}
func handleUserCommand(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if !strings.HasPrefix(m.Content, config.DefaultPrefix) {
		return false
	}
	if m.Content == config.DefaultPrefix {
		return false
	}
	m.Content = m.Content[len(config.DefaultPrefix):]
	args := strings.Split(m.Content, " ")
	if len(args) < 1 {
		return false
	}
	switch args[0] {
	case "remind":
		return handleRemindCommand(s, m, args[1:])
	default:
		return false
	}
}

func isPermitted(guildID, UserID string) bool {
	//TODO: make function use database for multiple guilds
	member, err := session.GuildMember(guildID, UserID)
	if err != nil {
		log.Println("isPermitted Unable to get guild member! ", err)
		return false
	}
	for _, role := range member.Roles {
		switch role {
		case "412193755286732800", "422408722107596811":
			return true
		}
	}
	return false
}

func handleRemindCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) bool {
	if !isPermitted(m.GuildID, m.Author.ID) {
		sendTimeoutMessage(m.ChannelID, "<:blobstop:657036452303208508>", 3*time.Second)
		return false
	}
	if len(args) < 4 {
		printHelp(m.ChannelID)
		return false
	}
	guild, err := s.Guild(m.GuildID)
	if err != nil {
		log.Println("handleRemindCommand Unable to get guild! ", err)
		return false
	}
	var message *discordgo.Message
	for _, channel := range guild.Channels {
		message, err = session.ChannelMessage(channel.ID, args[0])
		if err == nil {
			break
		}
	}
	if err != nil {
		log.Println("handleRemindCommand Unable to load message ID: "+args[0]+"! ", err)
		_, err = s.ChannelMessageSend(m.ChannelID, "Unable to find message with this ID :worried:")
		if err != nil {
			log.Println("handleRemindCommand Unable to send channel message! ", err)
		}
		return false
	}
	reactionId := args[1]
	if strings.HasPrefix(args[1], "<") {
		reactionId = args[1][:len(args[1])-1]
	}
	afterID := ""
	for {
		reactionists, err := session.MessageReactions(m.ChannelID, message.ID, reactionId, 100, "", afterID)
		if err != nil {
			log.Println("handleRemindCommand Unable to get message reactionists! ", err)
			return false
		}
		member, err := s.GuildMember(m.GuildID, m.Author.ID)
		if err != nil {
			log.Println("handleRemindCommand Unable to get guild member! ", err)
			return false
		}
		content := strings.Join(args[2:], " ")
		for _, reactionist := range reactionists {
			go sendReminder(member, guild, reactionist.ID, content)
		}
		if len(reactionists) == 100 {
			afterID = reactionists[99].ID
		} else {
			break
		}
	}
	return true
}

func sendReminder(remindingMember *discordgo.Member, guild *discordgo.Guild, userID, content string) bool {
	dm, err := session.UserChannelCreate(userID)
	if err != nil {
		log.Println("sendReminder Unable to create user channel! ", err)
		return false
	}
	nick := remindingMember.Nick
	if len(nick) == 0 {
		nick = remindingMember.User.Username
	}
	embed := discordgo.MessageEmbed{
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: guild.IconURL(),
		},
		Title:       nick + " reminds you on " + guild.Name + ":",
		Description: content,
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       0x00FFFF,
	}
	_, err = session.ChannelMessageSendEmbed(dm.ID, &embed)
	if err != nil {
		log.Println("sendReminder Unable to send channel message embed! ", err)
		return false
	}
	return true
}
