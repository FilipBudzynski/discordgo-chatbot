package commands

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type CommandID int

// Const to add different commands
const (
	PingCommandID CommandID = iota
	HelpCommandID
	PlayCommandID
	PauseCommandID
	ResumeCommandID
	UnknownCommandID
)

type Command struct {
	Message   *discordgo.MessageCreate
	Args      []string
	CommandID CommandID
}

// Store voice instances here
var (
	voiceInstances     = make(map[string]*VoiceInstance)
	voiceInstanceMutex sync.Mutex
)

func ParseCommand(content string) CommandID {
	switch content {
	case "!ping":
		return PingCommandID
	case "!help", "!h":
		return HelpCommandID
	case "!play", "!p":
		return PlayCommandID
	case "!pause":
		return PauseCommandID
	case "!resume":
		return ResumeCommandID
	default:
		return UnknownCommandID
	}
}

// Handles commands based on the commandID
func CommandHandler(s *discordgo.Session, commandChan <-chan Command) {
	for c := range commandChan {
		guildID := c.Message.GuildID
		authorID := c.Message.Author.ID

		// TODO: in voice related functions, check if the vs is not nil, if nil send message telling user to join the voice channel
		vs, err := s.State.VoiceState(guildID, authorID)
		if err != nil {
			fmt.Println("Could not find the VoiceState", err)
			return
		}

		switch c.CommandID {
		case PingCommandID:
			go func(channelID string) {
				go handlePingCommand(s, channelID)
			}(c.Message.ChannelID)
		case PlayCommandID:
			go handlePlayCommand(s, vs, guildID, authorID, c.Args[1])
		case PauseCommandID:
			vi := getVoiceInstancce(vs.ChannelID)
			go vi.Pause()
		case ResumeCommandID:
			vi := getVoiceInstancce(vs.ChannelID)
			go vi.Resume()
		default:
			err := sendUnknownCommand(s, c.Message.ChannelID)
			if err != nil {
				fmt.Println("Error with SendPong command", err)
			}
		}
	}
}

func getVoiceInstancce(voiceChannelID string) *VoiceInstance {
	voiceInstanceMutex.Lock()
	vi := voiceInstances[voiceChannelID]
	voiceInstanceMutex.Unlock()
	return vi
}
