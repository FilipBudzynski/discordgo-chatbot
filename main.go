package main

import (
	"discord_go_chat/commands"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	Token       string
	commandChan = make(chan commands.Command)
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()

	if Token == "" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file")
			os.Exit(1)
		}
		Token = os.Getenv("DISCORD_TOKEN")
	}
}

func main() {
	// Create a new discord session
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error createing Discord session", err)
		return
	}

	go commands.CommandHandler(dg, commandChan)

	dg.AddHandler(handleMessage)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	defer dg.Close()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

// Checks whether the written message is a command and parses it
// in order to create a Command struct which will be send to
// command channel to invoke actions
func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Do not parse bot messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if the message is a command
	if !strings.HasPrefix(m.Content, "!") {
		return
	}

	commandID := commands.ParseCommand(m.Content)
	command := commands.Command{
		CommandID: commandID,
		AuthorID:  m.Author.ID,
		GuildID:   m.GuildID,
		ChannelID: m.ChannelID,
	}

	commandChan <- command
}

func UserInVoiceChannel(s *discordgo.Session, guildID, userID string) (result bool, err error) {
	voiceState, err := s.State.VoiceState(guildID, userID)
	if err != nil {
		return false, err
	} else {

		if voiceState != nil {
			return true, nil
		}

		fmt.Println("User not in channel")
		return false, nil
	}
}
