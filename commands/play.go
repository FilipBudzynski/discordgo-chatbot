package commands

import (
	"discord_go_chat/audio"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/bwmarrin/discordgo"
)

// TODO: now the voice instance is repsponsible for both playing sounds and also generating url
// for the song. This needs to be changed to seperate resposibility by creating a song type structure that will
// have link and some other stuff

// Type used for creating a voice instance that is responsible for playing sounds/songs on the channel
type VoiceInstance struct {
	Session         *discordgo.Session
	VoiceConnection *discordgo.VoiceConnection
	VoiceState      *discordgo.VoiceState
	PlaybackState   *audio.PlaybackState
	GuildID         string
	VoiceChannelID  string
	AuthorID        string
	// Queue           chan string
}

func handlePlayCommand(s *discordgo.Session, vs *discordgo.VoiceState, guildID, authorID, ytLink string) {
	vi := VoiceInstance{
		Session:        s,
		VoiceState:     vs,
		GuildID:        guildID,
		VoiceChannelID: vs.ChannelID,
		AuthorID:       authorID,
	}

	// store voice instance
	voiceInstanceMutex.Lock()
	voiceInstances[vs.ChannelID] = &vi
	voiceInstanceMutex.Unlock()

	vi.PlayLink(ytLink)
}

// establishes voiceConnection and plays song from provided URL
func (v *VoiceInstance) PlayLink(youtubeURL string) {
	// TODO: check if bot needs to be connected to the channel or is already connected
	vc, err := v.Session.ChannelVoiceJoin(v.GuildID, v.VoiceChannelID, false, false)
	if err != nil {
		fmt.Println("Error joining voice channel:", err)
		return
	}

	v.VoiceConnection = vc

	url, _ := GetAudioURL(youtubeURL)
	err = v.sendOpusAudio(url)
	if err != nil {
		fmt.Println("Error with sending Opus Audio", err)
	}

	err = vc.Disconnect()
	if err != nil {
		fmt.Println("Error with disconecting from voice channel", err)
	}
}

// Function streams sound directly from the link
func (v *VoiceInstance) sendOpusAudio(url string) error {
	time.Sleep(250 * time.Millisecond)

	err := v.VoiceConnection.Speaking(true)
	if err != nil {
		log.Fatal("Faild setting speaking to true", err)
	}

	done := make(chan bool)
	state := audio.NewMutexState()
	v.PlaybackState = state

	fmt.Println("Audio is playing")
	audio.PlayAudioFile(v.VoiceConnection, url, done, state)

	select {
	case <-done:
		err = v.VoiceConnection.Speaking(false)
		if err != nil {
			log.Fatal("Faild setting speaking to false", err)
		}

		// wait before disconecting
		time.Sleep(250 * time.Millisecond)

		return nil
	}
}

func GetAudioURL(videoURL string) (string, error) {
	youtubeDownloader, err := exec.LookPath("yt-dlp")
	if err != nil {
		fmt.Println("yt-dlp not found in path")
		return "", err
	}

	args := []string{
		"--get-title",
		"--get-duration",
		"--get-thumbnail",
		"--extract-audio",
		"--audio-format", "best",
		"--get-url", videoURL,
	}

	cmd := exec.Command(youtubeDownloader, args...)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	fmt.Print(string(output))
	return string(output), nil
}

func (v *VoiceInstance) Pause() {
	v.PlaybackState.Pause()
}

func (v *VoiceInstance) Resume() {
	v.PlaybackState.Resume()
}
