package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Action given by the command "!ping", sends back a "Pong!" response
func SendPong(s *discordgo.Session, channelID string) error {
	_, err := s.ChannelMessageSend(channelID, "Pong!")
	if err != nil {
		fmt.Println("Error sending message to channel:", err)
		return err
	}
	return nil
}
