package castigate

type Config struct {
	Podcasts           []*Podcast
	FilenameTemplate   string
	DefaultCountToKeep int
}
