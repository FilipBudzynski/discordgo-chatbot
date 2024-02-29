package music

import (
	"fmt"
	"os/exec"
	"strings"
)

type Song struct {
	URL       string
	Title     string
	Thumbnail string
	Duration  string
}

func (s Song) String() string {
	return fmt.Sprintf("Title: %s\nURL: %s\nDuration: %s\nThumbnail: %s\n", s.Title, s.URL, s.Duration, s.Thumbnail)
}

func NewSong(songData string) Song {
	s := Song{}

	lines := strings.Split(songData, "\n")

	s.Title = lines[0]
	s.URL = lines[1]
	s.Thumbnail = lines[2]
	s.Duration = lines[3]

	return s
}

func GetSongInfo(videoURL string) (string, error) {
	ytdlp, err := exec.LookPath("yt-dlp")
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

	cmd := exec.Command(ytdlp, args...)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	fmt.Print(string(output))
	return string(output), nil
}
