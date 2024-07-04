package cmd

import (
	"castigate/feed"
	"os"
	"testing"
)

func TestEdit(t *testing.T) {
	fn, config := CreateTestConfigFile(t)
	defer os.Remove(fn)

	backend := feed.FileBackend{}
	backend.Init(fn)
	config, err := backend.Load()
	if err != nil {
		t.Fatal(err)
	}
	config.Podcasts = append(config.Podcasts, &feed.Podcast{
		Label:       "test",
		Feed:        "http://feed.example.com",
		Directory:   "test",
		CountToKeep: 0,
		Start:       "oldest",
		Episodes:    make(map[string]*feed.Episode, 0),
	})
	err = backend.Save(config)
	if err != nil {
		t.Fatal(err)
	}

	rootCmd.SetArgs([]string{"--config", fn, "edit", "test", "--url", "http://feed.example.com/rss", "--count", "42", "--directory", "foo",
		"--start", "newest"})
	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("error editing podcast test: %v", err)
	}

	config, err = backend.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(config.Podcasts) != 1 {
		t.Fatalf("expected 1 podcast, got %d", len(config.Podcasts))
	}
	podcast := config.Podcasts[0]
	if podcast.Feed != "http://feed.example.com/rss" {
		t.Fatalf("expected http://feed.example.com/rss feed, got %s", podcast.Feed)
	}
	if podcast.Start != "newest" {
		t.Fatalf("expected newest, got %s", podcast.Start)
	}
	if podcast.CountToKeep != 42 {
		t.Fatalf("expected 42 episodes, got %d", podcast.CountToKeep)
	}
	if podcast.Directory != "foo" {
		t.Fatalf("expected foo directory, got %s", podcast.Directory)
	}
}
