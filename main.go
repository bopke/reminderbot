package main

import (
	"encoding/json"
	"github.com/bopke/discordgo"
	"github.com/robfig/cron"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	MysqlString   string `json:"mysql_string"`
	DiscordToken  string `json:"discord_token"`
	DefaultPrefix string `json:"default_prefix"`
}

var config Config

var session *discordgo.Session

func loadConfig() (c Config) {
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Panic(err)
	}
	defer configFile.Close()
	err = json.NewDecoder(configFile).Decode(&c)
	if err != nil {
		log.Panic("loadConfig Unable to decode config! ", err)
	}
	return
}

func main() {
	var err error
	log.Println("Starting...")
	config = loadConfig()
	//InitDB()

	session, err = discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		log.Panicln(err)
	}

	session.AddHandler(OnMessageCreate)
	//session.AddHandler(OnGuildCreate)

	err = session.Open()
	if err != nil {
		log.Panicln(err)
	}

	c := cron.New()
	_ = c.AddFunc("* * * * *", check)
	c.Start()

	go check()

	log.Println("Started!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	log.Println("Caught stop signal")
	err = session.Close()
	if err != nil {
		log.Panicln(err)
	}
}

//TODO: Set more descriptive name for that function
func check() {

}

func printHelp(channelID string) {
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0xFF0000,
		Title:       "Reminder bot commands",
		Description: "!remind <messageID> <emoji> <reminder content>",
		Timestamp:   time.Now().Format(time.RFC3339),
	}
	_, err := session.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		log.Println("printHelp Unable to send embed! ", err)
	}
}

func sendTimeoutMessage(channelID, content string, duration time.Duration) {
	msg, err := session.ChannelMessageSend(channelID, content)
	if err != nil {
		log.Println("handleRemindCommand Unable to send channel message! ", err)
		return
	}
	time.Sleep(duration)
	err = session.ChannelMessageDelete(channelID, msg.ID)
	if err != nil {
		log.Println("handleUserCommand Unable to delete channel message! ", err)
	}
}
