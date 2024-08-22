package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	Api "github.com/Nvim/silverstalker/Api"
	"github.com/bwmarrin/discordgo"
)

var (
	BotToken string
	Bot      *discordgo.Session
)

func Init() (err error) {
	Bot, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal("Error creating bot: ", err.Error())
	}
	return err
}

func Listen() error {
	// create a session
	discord := Bot

	// add a event handler
	discord.AddHandler(newMessage)

	// open session
	err := discord.Open()
	if err != nil {
		log.Fatal("Error opening discord websocket")
		return err
	}
	defer discord.Close() // close session, after function termination

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	return nil
}

func SendMessage(msg string) error {
	err := Bot.Open()
	if err != nil {
		log.Fatal("Error opening discord websocket")
		return err
	}
	defer Bot.Close() // close session, after function termination
	_, err = Bot.ChannelMessageSend("1273632829753917515", msg)
	if err != nil {
		log.Fatal("couldn't send message in channel")
	}
	return nil
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// respond to user message if it contains `!help` or `!bye`
	switch {
	case strings.Contains(message.Content, "!help"):
		_, err := discord.ChannelMessageSend(message.ChannelID, "Hello World")
		if err != nil {
			log.Fatal("couldn't send message in channel")
		}
	case strings.Contains(message.Content, "!lucas"):
		s, err := Api.Api()
		if err != nil {
			fmt.Println("Error trace: " + err.Error())
		}
		_, err = discord.ChannelMessageSend(message.ChannelID, s)
		if err != nil {
			log.Fatal("couldn't send message in channel")
		}
	}
}
