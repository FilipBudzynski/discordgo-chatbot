package music

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Song struct {
	URL       string
	Title     string
	Thumbnail string
	Duration  string
	Path      string
}

func (s Song) String() string {
	return fmt.Sprintf("Title: %s\nURL: %s\nDuration: %s\nThumbnail: %s\nPath: %s\n", s.Title, s.URL, s.Duration, s.Thumbnail, s.Path)
}

func NewSong(songData, audioPath string) Song {
	s := Song{}

	lines := strings.Split(songData, "\n")

	s.Title = lines[0]
	s.URL = lines[1]
	s.Thumbnail = lines[2]
	s.Duration = lines[3]
	s.Path = audioPath

	return s
}

func GetSongData(videoURL string) (string, error) {
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

	// fmt.Print(string(output))
	return string(output), nil
}

func DownloadAudio(videoURL string) (string, error) {
	ytdlp, err := exec.LookPath("yt-dlp")
	if err != nil {
		fmt.Println("yt-dlp not found in path")
		return "", err
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return "", err
	}
	fileName := extractVideoID(videoURL)

	audioPath := currentDir + "/" + fileName + ".mp3"

	args := []string{
		"--extract-audio",
		"--audio-format", "mp3",
		"--output", audioPath,
		videoURL,
	}

	cmd := exec.Command(ytdlp, args...)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error with getting cmd output")
		return "", err
	}

	fmt.Println("yt-dlp output:", string(output))
	fmt.Println("Audio file saved successfully at:", audioPath)

	return audioPath, nil
}

func extractVideoID(url string) string {
	index := strings.Index(url, "v=")
	if index == -1 {
		return "" // Video ID not found
	}

	url = url[index+2:] // Remove characters before "v="
	index = strings.Index(url, "_")
	if index == -1 {
		return url // If underscore not found, return the remaining characters
	}

	return url[:index] // Return characters up to the first underscore
}
