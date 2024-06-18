package feed

type Config struct {
	Podcasts           []*Podcast
	FilenameTemplate   string
	DefaultCountToKeep int
}
