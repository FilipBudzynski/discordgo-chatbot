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
	QueueCommandID
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
	case "!queue", "!q":
		return QueueCommandID
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
				go HandlePingCommand(s, channelID)
			}(c.Message.ChannelID)
		case PlayCommandID:
			ytLink := c.Args[1]
			vi := getVoiceInstance(vs.ChannelID)
			if vi == nil {
				v := NewVoiceInstance(s, vs, guildID, authorID)

				voiceInstanceMutex.Lock()
				voiceInstances[vs.ChannelID] = v
				voiceInstanceMutex.Unlock()

				v.init()
				vi = v
			}

			go vi.play(ytLink)

		case PauseCommandID:
			vi := getVoiceInstance(vs.ChannelID)
			if vi == nil {
				fmt.Println("Voice instance not initiated")
			}
			go vi.Pause()
		case ResumeCommandID:
			vi := getVoiceInstance(vs.ChannelID)
			if vi == nil {
				fmt.Println("Voice instance not initiated")
			}
			go vi.Resume()
		case QueueCommandID:
			vi := getVoiceInstance(vs.ChannelID)
			if vi == nil {
				fmt.Println("Voice instance not initiated")
			}
			go vi.showQueue()

		default:
			err := sendUnknownCommand(s, c.Message.ChannelID)
			if err != nil {
				fmt.Println("Error with SendPong command", err)
			}
		}
	}
}

func getVoiceInstance(voiceChannelID string) *VoiceInstance {
	voiceInstanceMutex.Lock()
	vi := voiceInstances[voiceChannelID]
	voiceInstanceMutex.Unlock()
	return vi
}
