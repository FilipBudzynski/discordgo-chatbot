package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Util function to send message of certain content
func SendMessage(s *discordgo.Session, channelID, message string) error {
	_, err := s.ChannelMessageSend(channelID, message)
	if err != nil {
		fmt.Println("Error sending message to channel:", err)
		return err
	}
	return nil
}
