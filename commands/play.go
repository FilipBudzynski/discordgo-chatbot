package commands

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

var (
	filePath = "audio.opus"
	buffer   = make([][]byte, 0)
)

func PlayLink(s *discordgo.Session, youtubeLink, GuildID, AuthorID, VoiceChannelID string) {
	// ytID := extractID(youtubeLink)
	// videoURL := "https://www.youtube.com/watch?v=" + ytID

	// Run yt-dlp command to get direct audio stream link
	// downloadAudio(youtubeLink)

	// TODO: check if bot needs to be connected to the channel or is already connected
	// Open a connection to the voice channel
	vc, err := s.ChannelVoiceJoin(GuildID, VoiceChannelID, false, false)
	if err != nil {
		fmt.Println("Error joining voice channel:", err)
		return
	}

	// Send Opus audio stream to the voice channel
	url, _ := getAudioURL(youtubeLink)
	sendOpusAudio(vc, url)

	vc.Disconnect()
}

func extractID(url string) string {
	parts := strings.Split(url, "?v=")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

func downloadAudio(videoURL string) {
	youtubeDownloader, err := exec.LookPath("yt-dlp")
	if err != nil {
		fmt.Println("yt-dlp not found in path")
		return
	}
	// Download audio in the best available format
	args := []string{
		"--extract-audio",
		"--audio-format", "opus",
		"--output", "audio",
		"--ignore-errors",
		videoURL,
	}

	cmd := exec.Command(youtubeDownloader, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// // Get the downloaded audio file name
	// files, err := ioutil.ReadDir(".")
	// if err != nil {
	// 	fmt.Println("Error reading directory:", err)
	// 	return
	// }
	//
	// var audioFileName string
	// for _, file := range files {
	// 	if strings.HasPrefix(file.Name(), "audio.") {
	// 		audioFileName = file.Name()
	// 		break
	// 	}
	// }
	//
	// if audioFileName == "" {
	// 	fmt.Println("No audio file found")
	// 	return
	// }
	//
	// // Convert the downloaded audio to DCA format using ffmpeg
	// // ffmpegCmd := exec.Command("ffmpeg", "-i", audioFileName, "-f", "s16le", "-ar", "48000", "-ac", "2", "audio.dca")
	// ffmpegCmd := exec.Command("ffmpeg", "-i", audioFileName, "-f", "s16le", "-ar", "48000", "-ac", "2", "-acodec", "pcm_s16le", "-b:a", "128k", "audio.dca")
	// ffmpegCmd.Stdout = os.Stdout
	// ffmpegCmd.Stderr = os.Stderr
	//
	// // Execute ffmpeg command
	// err = ffmpegCmd.Run()
	// if err != nil {
	// 	fmt.Println("Error converting audio to DCA:", err)
	// 	return
	// }
	//
	// fmt.Println("Audio downloaded and converted to .dca format successfully.")
}

func loadSoundOpus(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening audio file:", err)
		return nil, err
	}
	defer file.Close()

	// Read the Opus audio data
	opusData, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading Opus data:", err)
		return nil, err
	}

	return opusData, nil
}

func loadSound(filePath string) error {
	file, err := os.Open("audio.opus")
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)
		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}

// func sendOpusAudio(vc *discordgo.VoiceConnection, filePath string) error {
// 	// Sleep for a specified amount of time before playing the sound
// 	time.Sleep(250 * time.Millisecond)
//
// 	// Start speaking.
// 	vc.Speaking(true)
//
// 	// Send the buffer data.
// 	for _, buff := range buffer {
// 		vc.OpusSend <- buff
// 	}
//
// 	// Stop speaking
// 	vc.Speaking(false)
//
// 	// Sleep for a specificed amount of time before ending.
// 	time.Sleep(250 * time.Millisecond)
//
// 	// Disconnect from the provided voice channel.
// 	vc.Disconnect()
//
// 	return nil
// }

func sendOpusAudio(vc *discordgo.VoiceConnection, url string) error {
	time.Sleep(250 * time.Millisecond)

	vc.Speaking(true)

	err := DCA(vc, url)
	if err != nil {
		log.Println("Failed to stream audio:", err)
	}

	vc.Speaking(false)

	time.Sleep(250 * time.Millisecond)

	return nil
}

func DCA(vc *discordgo.VoiceConnection, url string) error {
	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 96
	opts.Application = "lowdelay"

	encodeSession, err := dca.EncodeFile(url, opts)
	if err != nil {
		return fmt.Errorf("failed creating an encoding session: %v", err)
	}

	done := make(chan error)
	dca.NewStream(encodeSession, vc, done)

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
	cmd := exec.Command("yt-dlp", "--extract-audio", "--audio-format", "best", "--get-url", videoURL)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
