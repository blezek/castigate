/*
Copyright © 2023 Daniel Blezek <blezek.daniel@mayo.edu>
This file is part of a CLI application.
*/
package cmd

import (
	"castigate/feed"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove podcasts",
	Args:  cobra.MinimumNArgs(1),
	Run:   runRemoveCmd,
}

func runRemoveCmd(cmd *cobra.Command, args []string) {
	log.Debugf("run command remove")
	backend, config := LoadConfiguration(cmd)
	for _, label := range args {
		for index, podcast := range config.Podcasts {
			if podcast.Label == label {
				log.Debugf("removing podcast %s", label)
				// see https://stackoverflow.com/a/57213476
				p := make([]*feed.Podcast, 0)
				p = append(p, config.Podcasts[:index]...)
				config.Podcasts = append(p, config.Podcasts[index+1:]...)
				break
			}
		}
	}
	err := backend.Save(config)
	if err != nil {
		log.Fatalf("error saving config: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
