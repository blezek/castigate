package cmd

import (
	"bytes"
	"castigate/feed"
	"os"
	"testing"
)

func CreateTestConfigFile(t *testing.T) (string, feed.Config) {
	fn, err := os.CreateTemp("", "test_padcast_feed")
	if err != nil {
		t.Fatal(err)
	}
	fn.Close()
	backend := feed.FileBackend{}
	backend.Init(fn.Name())
	config := feed.NewConfig()
	backend.Save(config)
	return fn.Name(), config
}

func TestRemove(t *testing.T) {
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
	rootCmd.SetArgs([]string{"--config", fn, "remove", "test"})
	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("error removing podcast test: %v", err)
	}

	config, err = backend.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(config.Podcasts) != 0 {
		t.Fatalf("expected 0 podcast, got %d", len(config.Podcasts))
	}
}

func TestRemoveNoneExistant(t *testing.T) {
	fn, config := CreateTestConfigFile(t)
	defer os.Remove(fn)

	buffer := new(bytes.Buffer)
	rootCmd.SetOut(buffer)
	rootCmd.SetErr(buffer)
	rootCmd.SetArgs([]string{"--config", fn, "remove", "test"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("error removing podcast test: %v", err)
	}

	backend := feed.FileBackend{}
	backend.Init(fn)
	config, err = backend.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(config.Podcasts) != 0 {
		t.Fatalf("expected 0 podcast, got %d", len(config.Podcasts))
	}
}
