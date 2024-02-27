package song

import (
	"fmt"
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

func (s *Song) parseJson(buffer string) {
	// TODO: brute force for now

	lines := strings.Split(buffer, "\n")
	s.Title = lines[0]
	s.URL = lines[1]
	s.Thumbnail = lines[2]
	s.Duration = lines[3]
}
