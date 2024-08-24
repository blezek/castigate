package cmd

import (
	"castigate/feed"
	"errors"
	"fmt"
	"github.com/gorilla/feeds"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func GetRSS(url string, t *testing.T) string {

	rss := feeds.Feed{
		Title:       "test",
		Link:        &feeds.Link{Href: url + "/rss"},
		Description: "test",
		Author:      nil,
		Updated:     time.Time{},
		Created:     time.Time{},
		Id:          "1234",
		Subtitle:    "subtitle",
		Items:       make([]*feeds.Item, 0),
		Copyright:   "copyright me 2024",
		Image:       nil,
	}
	startDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	for count := 0; count < 100; count++ {
		episode := fmt.Sprintf("episode-%03d", count)
		itemUrl := url + fmt.Sprintf("/asset/%s.mp3", episode)
		link := &feeds.Link{Href: itemUrl}
		enclosure := &feeds.Enclosure{
			Url:    itemUrl,
			Length: "0",
			Type:   "audio/mpeg",
		}
		rssItem := &feeds.Item{
			Title:       episode + fmt.Sprintf(" this is episode #%d", count),
			Link:        link,
			Description: episode,
			Id:          episode,
			Updated:     startDate,
			Created:     startDate,
			Enclosure:   enclosure,
		}
		rss.Add(rssItem)
		startDate = startDate.AddDate(0, 0, 1)
	}
	rssText, err := rss.ToRss()
	if err != nil {
		t.Fatalf("convert feed to RSS: %v", err)
	}
	return rssText
}

func CreateTestServer(t *testing.T) *httptest.Server {
	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)
	mux.HandleFunc("/rss", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte(GetRSS(ts.URL, t)))
	})
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("asset"))
	})

	return ts
}

func TestParseRSS(t *testing.T) {
	ts := CreateTestServer(t)
	defer ts.Close()

	client := http.Client{}
	type testData struct {
		url      string
		expected string
	}
	tests := []testData{
		{ts.URL + "/rss", GetRSS(ts.URL, t)},
		{ts.URL + "/b/asset/foo.mp3", "asset"},
		{ts.URL + "/b/artwork/art.jpg", "asset"},
	}

	for _, test := range tests {

		request, err := http.NewRequest("GET", test.url, nil)
		if err != nil {
			t.Fatalf("could not get %s: %v", ts.URL, err)
		}
		response, err := client.Do(request)
		if err != nil {
			t.Fatalf("could not make request %s: %v", request.URL.String(), err)
		}
		if response.StatusCode != 200 {
			t.Fatalf("bad status code: %d", response.StatusCode)
		}
		b, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("could not read response body: %v", err)
		}
		if string(b) != test.expected {
			t.Fatalf("bad response body expected 'rss' got: %s", string(b))
		}
	}
}

func TestGetFeed(t *testing.T) {
	dir, err := os.MkdirTemp("", "test_padcast_feed")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	ts := CreateTestServer(t)
	defer ts.Close()
	config := feed.Config{
		Podcasts: []*feed.Podcast{
			{
				Label:       "test",
				Title:       "",
				Feed:        ts.URL + "/rss",
				Directory:   dir,
				CountToKeep: 0,
				Start:       "oldest",
				Episodes:    nil,
			},
		},
		FilenameTemplate:   feed.DefaultFilenameTemplate,
		DefaultCountToKeep: 10,
	}
	podcast := config.Podcasts[0]
	err = podcast.Sync(config, "")
	if err != nil {
		t.Fatalf("could not sync podcast: %v", err)
	}
	if podcast.GetNewCount() != 90 {
		t.Errorf("expected 90 got %d", podcast.GetNewCount())
	}
	if podcast.GetDownloadedCount() != 10 {
		t.Errorf("expected 10 got %d", podcast.GetDownloadedCount())
	}
	if podcast.GetDeletedCount() != 0 {
		t.Errorf("expected 0 got %d", podcast.GetDeletedCount())
	}
	filenames := []string{
		"2020-01-01-00-00-00-episode-000-this-is-episode--0.mp3",
		"2020-01-02-00-00-00-episode-001-this-is-episode--1.mp3",
		"2020-01-03-00-00-00-episode-002-this-is-episode--2.mp3",
		"2020-01-04-00-00-00-episode-003-this-is-episode--3.mp3",
		"2020-01-05-00-00-00-episode-004-this-is-episode--4.mp3",
		"2020-01-06-00-00-00-episode-005-this-is-episode--5.mp3",
		"2020-01-07-00-00-00-episode-006-this-is-episode--6.mp3",
		"2020-01-08-00-00-00-episode-007-this-is-episode--7.mp3",
		"2020-01-09-00-00-00-episode-008-this-is-episode--8.mp3",
		"2020-01-10-00-00-00-episode-009-this-is-episode--9.mp3",
	}
	for _, filename := range filenames {
		if !FileExists(filepath.Join(dir, filename)) {
			t.Errorf("missing file %s", filepath.Join(dir, filename))
		}
	}

	// sync again checking to see if we download more
	err = podcast.Sync(config, "")
	if err != nil {
		t.Fatalf("could not sync podcast: %v", err)
	}
	for _, filename := range filenames {
		if !FileExists(filepath.Join(dir, filename)) {
			t.Errorf("missing file %s", filepath.Join(dir, filename))
		}
	}
	downloadedFiles, err := filepath.Glob(dir + "/*.mp3")
	if err != nil {
		t.Errorf("can not find files is %s: %v", dir, err)
	}
	if len(downloadedFiles) != 10 {
		t.Errorf("expected 10 got %d downloaded files in %s", len(downloadedFiles), dir)
	}
	// delete a few and retry
	for _, filename := range filenames[:4] {
		os.Remove(filepath.Join(dir, filename))
	}

	_, err = podcast.UpdateFromRSS(config)
	if err != nil {
		t.Fatalf("could not sync podcast: %v", err)
	}
	if podcast.GetNewCount() != 90 {
		t.Errorf("expected 90 got %d", podcast.GetNewCount())
	}
	if podcast.GetDownloadedCount() != 10 {
		t.Errorf("expected 10 got %d", podcast.GetDownloadedCount())
	}
	if podcast.GetDeletedCount() != 0 {
		t.Errorf("expected 0 got %d", podcast.GetDeletedCount())
	}

	err = podcast.Sync(config, "")
	if err != nil {
		t.Fatalf("could not sync podcast: %v", err)
	}
	if podcast.GetNewCount() != 86 {
		t.Errorf("expected 86 got %d", podcast.GetNewCount())
	}
	if podcast.GetDownloadedCount() != 10 {
		t.Errorf("expected 10 got %d", podcast.GetDownloadedCount())
	}
	if podcast.GetDeletedCount() != 4 {
		t.Errorf("expected 4 deleted got %d", podcast.GetDeletedCount())
	}
}

func FileExists(fn string) bool {
	file, err := os.Open(fn)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	file.Close()
	return true

}
