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
	SkipCommandID
	QueueCommandID
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
	case "!queue", "!q":
		return QueueCommandID
	case "!skip", "!s":
		return SkipCommandID
	default:
		return UnknownCommandID
	}
}

// Handles commands based on the commandID
func CommandHandler(s *discordgo.Session, commandChan <-chan Command) {
	for c := range commandChan {
		guildID := c.Message.GuildID
		authorID := c.Message.Author.ID
		channelID := c.Message.ChannelID

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
				v := NewVoiceInstance(s, vs, guildID, authorID, channelID)

				fmt.Println("Createing new voice instance: ", len(voiceInstances))

				voiceInstanceMutex.Lock()
				voiceInstances[vs.ChannelID] = v
				voiceInstanceMutex.Unlock()

				go v.processQueue()
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
				err := SendMessage(s, channelID, "Connect to the voice channel to run this command.")
				if err != nil {
					fmt.Println("Error while sending message to the channel", err)
				}
				return
			}
			message := vi.printQueue()
			err := SendMessage(vi.Session, vi.ChannelID, message)
			if err != nil {
				fmt.Println("Couldn't send a message")
			}

		case SkipCommandID:
			go func() {
				vi := getVoiceInstance(vs.ChannelID)
				msg := fmt.Sprintf("Skipping: %s", vi.Queue[0].Title)
				err := SendMessage(vi.Session, vi.ChannelID, msg)
				if err != nil {
					fmt.Println("Couldn't send a message")
				}
				vi.SkipSong()
				return
			}()

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
