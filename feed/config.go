package feed

import "fmt"

type Config struct {
	Podcasts           []*Podcast
	FilenameTemplate   string
	DefaultCountToKeep int
}

func (c Config) FindPodcast(label string) (*Podcast, error) {
	for _, p := range c.Podcasts {
		if p.Label == label {
			return p, nil
		}
	}
	return nil, fmt.Errorf("podcast not found")
}
