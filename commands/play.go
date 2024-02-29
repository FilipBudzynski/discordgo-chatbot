package commands

import (
	"discord_go_chat/audio"
	"discord_go_chat/music"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var timeOut = 1 * time.Minute

// Type used for creating a voice instance that is responsible for playing sounds/songs on the channel
type VoiceInstance struct {
	Session         *discordgo.Session
	VoiceConnection *discordgo.VoiceConnection
	VoiceState      *discordgo.VoiceState
	Queue           chan music.Song
	PlaybackState   *audio.PlaybackState
	TimeoutDuration time.Duration
	IsPlaying       bool
	TimeoutSignal   chan bool
	Stop            chan bool
	GuildID         string
	VoiceChannelID  string
	AuthorID        string
}

func NewVoiceInstance(s *discordgo.Session, vs *discordgo.VoiceState, guildID, authorID string) *VoiceInstance {
	vi := &VoiceInstance{
		Session:         s,
		VoiceState:      vs,
		GuildID:         guildID,
		VoiceChannelID:  vs.ChannelID,
		AuthorID:        authorID,
		TimeoutDuration: timeOut,
		Queue:           make(chan music.Song),
	}
	return vi
}

func (v *VoiceInstance) init() {
	go v.processQueue()
}

// Establishes voiceConnection and plays song from provided URL
func (v *VoiceInstance) play(youtubeURL string) {
	v.joinVoiceChannel()

	songData, err := music.GetSongInfo(youtubeURL)
	if err != nil {
		fmt.Println("Error with getting info from yt-dlp: ", err)
	}

	song := music.NewSong(songData)
	v.Queue <- song
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
		case song, ok := <-v.Queue:
			if !ok {
				return
			}
			fmt.Println("Processing song:", song)
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(v.TimeoutDuration)

			v.playSong(song)
		}
	}
}

func (v *VoiceInstance) Disconnect() {
	v.VoiceConnection.Disconnect()

	voiceInstanceMutex.Lock()
	defer voiceInstanceMutex.Unlock()

	delete(voiceInstances, v.VoiceConnection.ChannelID)
}

func (v *VoiceInstance) playSong(s music.Song) {
	v.IsPlaying = true

	v.PlaybackState = audio.NewMutexState()
	audio.PlayAudioFile(v.VoiceConnection, s.URL, v.Stop, v.PlaybackState)

	v.IsPlaying = false
}

func (v *VoiceInstance) StopPlaying() {
	v.Stop <- true
}

func (v *VoiceInstance) Pause() {
	v.IsPlaying = false
	v.VoiceConnection.Speaking(false)
	v.PlaybackState.Pause()
}

func (v *VoiceInstance) Resume() {
	v.IsPlaying = true
	v.VoiceConnection.Speaking(true)
	v.PlaybackState.Resume()
}

func (v *VoiceInstance) showQueue() {
}
