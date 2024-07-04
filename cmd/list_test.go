package cmd

import (
	"bytes"
	"castigate/feed"
	"os"
	"testing"
)

func TestList(t *testing.T) {
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
	rootCmd.SetArgs([]string{"--config", fn, "list"})
	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("error listinf podcasts: %v", err)
	}

	podcast := config.Podcasts[0]
	s := podcast.PrintDetails()
	if s != expectedOutput {
		t.Fatalf("did not get expected output\nexpected:\n%s\nactual:\n%s", expectedOutput, s)
	}
}

const expectedOutput = `Title: 
Label: test
Feed: http://feed.example.com
Direction: oldest
Number of Episodes: 0
	Downloaded: 0
	New: 0
	Deleted: 0

`
