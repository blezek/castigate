package cmd

import (
	"bytes"
	"castigate/feed"
	"os"
	"testing"
)

func TestAddExisting(t *testing.T) {
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

	buffer := new(bytes.Buffer)
	rootCmd.SetOut(buffer)
	rootCmd.SetErr(buffer)
	rootCmd.SetArgs([]string{"--config", fn, "add", "test", "http://feed.example.com"})
	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("error adding podcast test: %v", err)
	}

	config, err = backend.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(config.Podcasts) != 1 {
		t.Fatalf("expected 1 podcast, got %d", len(config.Podcasts))
	}
}

func TestAddAdditional(t *testing.T) {
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

	buffer := new(bytes.Buffer)
	rootCmd.SetOut(buffer)
	rootCmd.SetErr(buffer)
	rootCmd.SetArgs([]string{"--config", fn, "add", "new_podcast", "http://feed.example.com"})
	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("error adding podcast test: %v", err)
	}

	config, err = backend.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(config.Podcasts) != 2 {
		t.Fatalf("expected 2 podcasts, got %d", len(config.Podcasts))
	}
}

func TestAddNew(t *testing.T) {
	fn, config := CreateTestConfigFile(t)
	defer os.Remove(fn)

	backend := feed.FileBackend{}
	backend.Init(fn)
	config, err := backend.Load()
	if err != nil {
		t.Fatal(err)
	}
	err = backend.Save(config)
	if err != nil {
		t.Fatal(err)
	}

	buffer := new(bytes.Buffer)
	rootCmd.SetOut(buffer)
	rootCmd.SetErr(buffer)
	rootCmd.SetArgs([]string{"--config", fn, "add", "--count", "2", "test", "http://feed.example.com"})
	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("error adding podcast test: %v", err)
	}

	config, err = backend.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(config.Podcasts) != 1 {
		t.Fatalf("expected 1 podcast, got %d", len(config.Podcasts))
	}
	podcast := config.Podcasts[0]
	if podcast.Label != "test" {
		t.Fatalf("expected podcast label to be test, got %s", podcast.Label)
	}
	if podcast.Feed != "http://feed.example.com" {
		t.Fatalf("expected podcast feed to be http://feed.example.com, got %s", podcast.Feed)
	}
	if podcast.Directory != "test" {
		t.Fatalf("expected podcast directory to be test, got %s", podcast.Directory)
	}
	if podcast.CountToKeep != 2 {
		t.Fatalf("expected podcast count to be 2, got %d", podcast.CountToKeep)
	}
}
