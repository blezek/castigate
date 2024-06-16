package castigate

import (
	"bytes"
	"errors"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"sort"
	"text/template"
	"time"
)

type Podcast struct {
	Feed        string
	Directory   string
	CountToKeep int
	Start       string // oldest or newest
	Episodes    map[string]*Episode
}

func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true

}

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}
func (podcast *Podcast) Sync(config Config) error {
	if podcast.Episodes == nil {
		podcast.Episodes = make(map[string]*Episode)
	}
	if podcast.Start == "" {
		podcast.Start = "oldest"
	}
	// Load the podcast, figure out what's going on
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(podcast.Feed)
	if err != nil {
		log.Errorf("could not parse feed: %s", podcast.Feed)
	}

	// Update any new episodes
	for _, item := range feed.Items {
		// do we have the episode?
		if podcast.Episodes[item.GUID] == nil {
			// construct the episode
			var audioFileURL string
			for _, item := range item.Enclosures {
				// TODO: Need to make sure it's an audio link
				audioFileURL = item.URL
			}
			t, _ := time.Parse(time.RFC1123Z, item.Published)
			episode := &Episode{
				GUID:     item.GUID,
				URL:      audioFileURL,
				State:    New,
				Filename: "",
				Date:     t,
			}
			episode.Filename = podcast.FormatFilename(config.FilenameTemplate, episode, item)
			podcast.Episodes[item.GUID] = episode
		}
	}
	// Update any downloaded -> deleted
	countOfExistingFiles := 0

	for _, episode := range podcast.Episodes {
		fn := path.Join(podcast.Directory, episode.Filename)
		if episode.State == Downloaded && !IsFileExist(fn) {
			episode.State = Deleted
		}
		if episode.State == Downloaded {
			countOfExistingFiles++
		}
	}

	// Sort, and download whatever we need
	orderedEpisodes := make([]*Episode, 0, len(podcast.Episodes))
	for _, episode := range podcast.Episodes {
		orderedEpisodes = append(orderedEpisodes, episode)
	}
	sort.Slice(orderedEpisodes, func(a, b int) bool {
		after := orderedEpisodes[a].Date.After(orderedEpisodes[b].Date)
		if podcast.Start == "oldest" {
			return !after
		} else {
			return after
		}
	})

	// loop through and download what we can
	countToDownload := config.DefaultCountToKeep - countOfExistingFiles
	for _, episode := range orderedEpisodes {
		if episode.State == New && countToDownload > 0 {
			err = episode.Download(path.Join(podcast.Directory, episode.Filename))
			if err == nil {
				episode.State = Downloaded
				countToDownload--
			} else {
				log.Errorf("could not download episode %s from %s: %s", episode.Filename, episode.URL, err)
			}
		}
	}

	return nil
}

func (podcast *Podcast) FormatFilename(filenameTemplate string, episode *Episode, item *gofeed.Item) string {
	tmpl, err := template.New("filenameTemplate").Parse(filenameTemplate)
	if err != nil {
		log.Errorf("could not parse filename template: %s", filenameTemplate)
	}
	context := map[string]interface{}{
		"item":    item,
		"episode": episode,
		"podcast": podcast,
	}
	buffer := bytes.Buffer{}
	err = tmpl.Execute(&buffer, context)
	if err != nil {
		log.Errorf("could not execute filename template: %s", err)
	}
	return buffer.String()
}
