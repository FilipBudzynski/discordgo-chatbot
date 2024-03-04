package commands

import (
	"discord_go_chat/audio"
	"discord_go_chat/music"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var timeOut = 1 * time.Minute

// Type used for creating a voice instance that is responsible for playing sounds/songs on the channel
type VoiceInstance struct {
	Session         *discordgo.Session
	VoiceConnection *discordgo.VoiceConnection
	VoiceState      *discordgo.VoiceState
	PlaybackState   *audio.PlaybackState
	GuildID         string
	VoiceChannelID  string
	ChannelID       string
	AuthorID        string
	NewSongSignal   chan struct{}
	Stop            chan bool
	Queue           []music.Song
	TimeoutDuration time.Duration
}

func NewVoiceInstance(s *discordgo.Session, vs *discordgo.VoiceState, guildID, authorID, channelID string) *VoiceInstance {
	vi := &VoiceInstance{
		Session:         s,
		VoiceState:      vs,
		GuildID:         guildID,
		ChannelID:       channelID,
		AuthorID:        authorID,
		VoiceChannelID:  vs.ChannelID,
		TimeoutDuration: timeOut,
		Stop:            make(chan bool),
		NewSongSignal:   make(chan struct{}),
		Queue:           make([]music.Song, 0),
	}

	vi.PlaybackState = audio.NewMutexState()
	return vi
}

func (v *VoiceInstance) init() {
	go v.processQueue()
}

func (v *VoiceInstance) play(youtubeURL string) {
	// connect if not yet connected to the channel
	if v.VoiceConnection == nil {
		v.joinVoiceChannel()
	}

	audioPath, err := music.DownloadAudio(youtubeURL)
	if err != nil {
		fmt.Println("Error with getting info from yt-dlp: ", err)
	}

	songData, err := music.GetSongData(youtubeURL)
	if err != nil {
		fmt.Println("Error with getting info from yt-dlp: ", err)
	}

	v.Queue = append(v.Queue, music.NewSong(songData, audioPath))
	v.NewSongSignal <- struct{}{}
}

func (v *VoiceInstance) joinVoiceChannel() {
	vc, err := v.Session.ChannelVoiceJoin(v.GuildID, v.VoiceChannelID, false, false)
	if err != nil {
		fmt.Println("Error joining voice channel:", err)
	}

	v.VoiceConnection = vc
}

func (v *VoiceInstance) processQueue() {
	timer := time.NewTimer(v.TimeoutDuration)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			fmt.Println("Timer expired, disconnecting...")
			v.Disconnect()
			return
		case <-v.NewSongSignal:
			timer.Stop()

			song := &v.Queue[0]

			audio.PlayAudioFile(v.VoiceConnection, song.Path, v.Stop, v.PlaybackState)

			// remove finished audio
			err := os.Remove(song.Path)
			if err != nil {
				fmt.Printf("Error while cleaning up: %v", err)
			}

			v.Queue = v.Queue[1:]

			timer.Reset(v.TimeoutDuration)
		}
	}
}

func (v *VoiceInstance) SkipSong() {
	v.Stop <- true
}

func (v *VoiceInstance) Pause() {
	v.VoiceConnection.Speaking(false)
	v.PlaybackState.Pause()
}

func (v *VoiceInstance) Resume() {
	v.VoiceConnection.Speaking(true)
	v.PlaybackState.Resume()
}

func (v *VoiceInstance) printQueue() string {
	if len(v.Queue) == 0 {
		fmt.Println("empty")
		return "Queue is empty, add music to queue with !play"
	}

	var message strings.Builder
	message.WriteString(fmt.Sprintln("Current Queue:"))

	for i, s := range v.Queue {
		message.WriteString(fmt.Sprintf("%d: %q\n", i+1, s.Title))
	}

	return message.String()
}

func (v *VoiceInstance) Disconnect() {
	err := v.VoiceConnection.Disconnect()
	if err != nil {
		fmt.Println("Erorr with disconnecting from the voice channel", err)
	}

	v.VoiceConnection.Close()

	voiceInstanceMutex.Lock()
	defer voiceInstanceMutex.Unlock()

	delete(voiceInstances, v.VoiceConnection.ChannelID)
}
