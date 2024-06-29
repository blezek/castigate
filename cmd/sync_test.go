package cmd

import (
	"castigate/feed"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func GetRSS(url string, t *testing.T) string {
	wd, _ := os.Getwd()
	log.Infof("in directory %s", wd)
	buffer, err := os.ReadFile("./testdata/test.rss")
	if err != nil {
		t.Fatalf("could not read test.rss: %v", err)
	}
	s := string(buffer)
	return strings.ReplaceAll(s, ":url:", url)
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
		FilenameTemplate:   DefaultFilenameTemplate,
		DefaultCountToKeep: 10,
	}
	podcast := config.Podcasts[0]
	err = podcast.Sync(config, "")
	if err != nil {
		t.Fatalf("could not sync podcast: %v", err)
	}
	if podcast.GetNewCount() != 118 {
		t.Errorf("expected 118 got %d", podcast.GetNewCount())
	}
	if podcast.GetDownloadedCount() != 10 {
		t.Errorf("expected 10 got %d", podcast.GetDownloadedCount())
	}
	if podcast.GetDeletedCount() != 0 {
		t.Errorf("expected 0 got %d", podcast.GetDeletedCount())
	}
	filenames := []string{
		"2022-01-18-05-01-00-124.-Biblical-Theology.mp3",
		"2022-01-25-05-01-00-21.-The-Intermediate-State.mp3",
		"2022-02-01-05-01-00-125.-Theophany-and-Christophany.mp3",
		"2022-02-08-05-01-00-22.-Amillennialism.mp3",
		"2022-02-15-05-01-00-126.-Adonai.mp3",
		"2022-02-22-05-01-00-23.-Assurance-of-Salvation.mp3",
		"2022-03-01-05-01-00-127.-Preterism.mp3",
		"2022-03-08-05-01-00-24.-Angels.mp3",
		"2022-03-15-04-01-00-128.-Universalism-and-Hell.mp3",
		"2022-03-22-04-01-00-25.-Repentance.mp3",
	}
	for _, filename := range filenames {
		if !FileExists(filepath.Join(dir, filename)) {
			t.Errorf("missing file %s", filepath.Join(dir, filename))
		}
	}

	// delete a few and retry
	for _, filename := range filenames[:4] {
		os.Remove(filepath.Join(dir, filename))
	}

	_, err = podcast.UpdateFromRSS(config)
	if err != nil {
		t.Fatalf("could not sync podcast: %v", err)
	}
	if podcast.GetNewCount() != 118 {
		t.Errorf("expected 118 got %d", podcast.GetNewCount())
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
	if podcast.GetNewCount() != 108 {
		t.Errorf("expected 108 got %d", podcast.GetNewCount())
	}
	if podcast.GetDownloadedCount() != 10 {
		t.Errorf("expected 10 got %d", podcast.GetDownloadedCount())
	}
	if podcast.GetDeletedCount() != 10 {
		t.Errorf("expected 10 got %d", podcast.GetDeletedCount())
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
