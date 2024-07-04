package cmd

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestRootCommand(t *testing.T) {
	buffer := new(bytes.Buffer)
	rootCmd.SetOut(buffer)
	rootCmd.SetErr(buffer)
	rootCmd.SetArgs([]string{"--debug", "help"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("Root command executed with error: %v", err)
	}
	if log.GetLevel() != log.DebugLevel {
		t.Errorf("Root command must set debug log level but was %v", log.GetLevel())
	}
}
