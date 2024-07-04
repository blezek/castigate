package feed

import (
	"fmt"
)

const DefaultFilenameTemplate = `{{.episode.Date.Format "2006-01-02-15-04-05" }}-{{.episode.Title}}.mp3`

type Config struct {
	Podcasts           []*Podcast
	FilenameTemplate   string
	DefaultCountToKeep int
}

func NewConfig() Config {
	return Config{
		Podcasts:           nil,
		FilenameTemplate:   DefaultFilenameTemplate,
		DefaultCountToKeep: 10,
	}
}
func (c Config) FindPodcast(label string) (*Podcast, error) {
	for _, p := range c.Podcasts {
		if p.Label == label {
			return p, nil
		}
	}
	return nil, fmt.Errorf("podcast not found")
}
