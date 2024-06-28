package feed

import (
	"fmt"
	"github.com/avast/retry-go/v4"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type EpisodeState int64

const (
	New EpisodeState = iota
	Downloaded
	Deleted
)

type Episode struct {
	GUID         string
	URL          string
	State        EpisodeState
	Title        string
	Filename     string
	Date         time.Time
	PodcastLabel string
}

func (episode Episode) String() string {
	return fmt.Sprintf("GUID: %s\nURL: %s\nState: %v\nFilename: %s\nDate: %s",
		episode.GUID, episode.URL, episode.State, episode.Filename, episode.Date)
}

func (episode *Episode) Download(path string) error {
	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0755)

	log.Debugf("Downloading %s to %s from %s", episode.Filename, dir, episode.URL)
	err := retry.Do(
		func() error {
			resp, err := http.Get(episode.URL)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			file, err := os.Create(path)
			if err != nil {
				return err
			}
			defer file.Close()
			count, err := io.Copy(file, resp.Body)
			log.Debugf("Downloaded %s to %s size %d", episode.Filename, path, count)
			return err
		})

	return err
}

/*
// Get data from Feed and put it into an Episode struct,
// returns a slice of episodes.
func GetEpisodes(feed *gofeed.Feed) []Episode {
	items := feed.Items

	episodes := []Episode{}

	for _, entry := range items {
		var audioFileURL string
		for _, item := range entry.Enclosures {
			// TODO: Need to make sure it's an audio link
			audioFileURL = item.URL
		}
		episode := Episode{
			PodcastTitle: feed.Title,
			//	Number:       counter,
			Title:   entry.Title,
			FileURL: audioFileURL,
			Date:    entry.Published,
		}
		episodes = append(episodes, episode)
	}

	// Sort slice by date in descending order
	sort.Slice(episodes, func(a, b int) bool {
		time1, _ := time.Parse(time.RFC1123Z, episodes[a].Date)
		time2, _ := time.Parse(time.RFC1123Z, episodes[b].Date)
		return time1.After(time2)
	})

	counter := len(episodes)
	for index := range episodes {
		counter--
		episodes[index].Number = counter
	}

	return episodes
}
*/
