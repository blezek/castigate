package feed

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"text/template"
	"time"
)

type Podcast struct {
	Label       string
	Title       string
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

func (podcast *Podcast) Sync(config Config, configFilePath string) error {
	if podcast.Episodes == nil {
		podcast.Episodes = make(map[string]*Episode, 0)
	}
	if podcast.Label == "" {
		podcast.Label = path.Base(podcast.Directory)
	}
	if podcast.Start == "" {
		podcast.Start = "oldest"
	}
	feed, err := podcast.UpdateFromRSS(config)
	if err != nil {
		return err
	}

	podcastDirectory := podcast.Directory
	if filepath.IsLocal(podcastDirectory) {
		absPath, err := filepath.Abs(configFilePath)
		if err != nil {
			log.Fatalf("could not get absolute path of %s", configFilePath)
		}
		podcastDirectory = filepath.Join(absPath, podcastDirectory)
	}

	// Update any downloaded -> deleted
	countOfExistingFiles := podcast.GetExistingFiles(podcastDirectory)

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

	log.Debugf("podcast directory is %s", podcastDirectory)
	// loop through and download what we can
	countToKeep := config.DefaultCountToKeep
	if podcast.CountToKeep > 0 {
		countToKeep = podcast.CountToKeep
	}
	re := regexp.MustCompile(`[^A-Za-z0-9_\-\.]`)
	countToDownload := countToKeep - countOfExistingFiles
	log.Infof("downloading %d episodes", countToDownload)
	for _, episode := range orderedEpisodes {
		if episode.State == New && countToDownload > 0 {
			path.Base(configFilePath)
			fn := re.ReplaceAllString(episode.Filename, "-")
			log.Infof("downloading %s to %s from %s", episode.Filename, fn, episode.URL)
			err = episode.Download(path.Join(podcastDirectory, fn))
			if err == nil {
				episode.State = Downloaded
				countToDownload--
			} else {
				log.Errorf("could not download episode %s from %s: %s", episode.Filename, episode.URL, err)
			}
		}
	}
	// save an m3u file
	playlistFilename := fmt.Sprintf("%s.m3u", feed.Title)
	playlistFilename = re.ReplaceAllString(playlistFilename, "-")
	fid, err := os.Create(path.Join(podcast.Directory, playlistFilename))
	if err != nil {
		fmt.Errorf("could not create the playlist: %s", err)
	}
	defer fid.Close()
	for _, episode := range orderedEpisodes {
		if episode.State == Downloaded {
			_, err = fid.WriteString(episode.Filename + "\n")
		}
	}

	return nil
}

func (podcast *Podcast) GetExistingFiles(podcastDirectory string) int {
	countOfExistingFiles := 0
	for _, episode := range podcast.Episodes {
		fn := path.Join(podcastDirectory, episode.Filename)
		if episode.State == Downloaded && !IsFileExist(fn) {
			episode.State = Deleted
		}
		if episode.State == Downloaded {
			countOfExistingFiles++
		}
	}
	return countOfExistingFiles
}

func (podcast *Podcast) UpdateFromRSS(config Config) (*gofeed.Feed, error) {
	log.Infof("fetching feed from %s", podcast.Feed)
	// Load the podcast, figure out what's going on
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(podcast.Feed)
	if err != nil {
		log.Errorf("could not parse feed: %s", podcast.Feed)
		log.Errorf("skipping podcast '%s'", podcast.Label)
		return nil, err
	}
	log.Infof("synchronizing %s", feed.Title)
	podcast.Title = feed.Title

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
				Title:    item.Title,
				Filename: "",
				Date:     t,
			}
			episode.Filename = podcast.FormatFilename(config.FilenameTemplate, episode, item)
			podcast.Episodes[item.GUID] = episode
		}
	}
	return feed, nil
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

func (podcast *Podcast) GetNewCount() int {
	counter := 0
	for _, episode := range podcast.Episodes {
		if episode.State == New {
			counter++
		}
	}
	return counter
}

func (podcast *Podcast) GetDownloadedCount() int {
	counter := 0
	for _, episode := range podcast.Episodes {
		if episode.State == Downloaded {
			counter++
		}
	}
	return counter

}

func (podcast *Podcast) GetDeletedCount() int {
	counter := 0
	for _, episode := range podcast.Episodes {
		if episode.State == Deleted {
			counter++
		}
	}
	return counter

}
