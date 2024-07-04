package cmd

import (
	"castigate/feed"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	fn, _ := CreateTestConfigFile(t)
	os.Remove(fn)

	rootCmd.SetArgs([]string{"--config", fn, "init", "--count", "3", "--template", "garf.mp3"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("error adding podcast test: %v", err)
	}

	backend := feed.FileBackend{}
	backend.Init(fn)
	config, err := backend.Load()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(fn)
	if config.FilenameTemplate != "garf.mp3" {
		t.Fatalf("Template filename should be garf.mp3 but got %s", config.FilenameTemplate)
	}
	if config.DefaultCountToKeep != 3 {
		t.Fatalf("Default count should be 3 but got %d", config.DefaultCountToKeep)
	}
}
