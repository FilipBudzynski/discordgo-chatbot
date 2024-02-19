package commands

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

// TODO: now the voice instance is repsponsible for both playing sounds and also generating url
// for the song. This needs to be changed to seperate resposibility by creating a song type structure that will
// have link and some other stuff

// Type used for creating a voice instance that is responsible for playing sounds/songs on the channel
type VoiceInstance struct {
	Session         *discordgo.Session
	VoiceConnection *discordgo.VoiceConnection
	VoiceState      *discordgo.VoiceState
	GuildID         string
	VoiceChannelID  string
	AuthorID        string
}

// Connects to voiceChannel specified in VoiceInstance,
// establishes voiceConnection and plays song from provided URL
func (v *VoiceInstance) PlayLink(youtubeURL string) {
	// TODO: check if bot needs to be connected to the channel or is already connected
	vc, err := v.Session.ChannelVoiceJoin(v.GuildID, v.VoiceChannelID, false, false)
	if err != nil {
		fmt.Println("Error joining voice channel:", err)
		return
	}

	v.VoiceConnection = vc

	url, _ := getAudioURL(youtubeURL)
	v.sendOpusAudio(url)

	vc.Disconnect()
}

// Function streams sound directly from the link
func (v *VoiceInstance) sendOpusAudio(url string) error {
	time.Sleep(250 * time.Millisecond)

	v.VoiceConnection.Speaking(true)

	err := v.streamDCA(url)
	if err != nil {
		log.Println("Failed to stream audio:", err)
	}

	v.VoiceConnection.Speaking(false)

	time.Sleep(250 * time.Millisecond)

	return nil
}

func (v *VoiceInstance) streamDCA(url string) error {
	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 96
	opts.Application = "lowdelay"

	encodeSession, err := dca.EncodeFile(url, opts)
	if err != nil {
		return fmt.Errorf("failed creating an encoding session: %v", err)
	}

	done := make(chan error)
	dca.NewStream(encodeSession, v.VoiceConnection, done)

	for {
		select {
		case err := <-done:
			if err != nil && err != io.EOF {
				log.Println("FATA: An error occured", err)
			}
			// Clean up incase something happened and ffmpeg is still running
			encodeSession.Cleanup()
			return nil
		}
	}
}

func getAudioURL(videoURL string) (string, error) {
	// youtubeDownloader, err := exec.LookPath("yt-dlp")
	// if err != nil {
	// 	fmt.Println("yt-dlp not found in path")
	// 	return "", err
	// }
	//
	// args := []string{
	// 	"--extract-audio",
	// 	"--audio-format", "best",
	// 	"--get-url", videoURL,
	// }
	//
	// cmd := exec.Command(youtubeDownloader, args...)
	cmd := exec.Command("yt-dlp", "--extract-audio", "--audio-format", "best", "--get-url", videoURL)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
