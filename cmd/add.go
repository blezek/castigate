/*
Copyright © 2023 Daniel Blezek <blezek.daniel@mayo.edu>
This file is part of a CLI application.
*/
package cmd

import (
	"castigate/feed"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a podcast to the configuration",
	Long: `Add a podcast to the configuration.
            arguments are <label> <url> [directory].  The label must be unique.
            if [directory] is not set, it defaults to label
              --count is the number of episodes to keep on disk, defaults to config if 0
              --direction is "oldest" or "newest" and dictates the order of episodes to download
            
            example:
               castigate add 5_minutes https://5minutesinchurchhistory.ligonier.org/rss`,
	Args: cobra.RangeArgs(2, 3),
	Run:  runAddCmd,
}

func runAddCmd(cmd *cobra.Command, args []string) {
	backend, config := LoadConfiguration(cmd)
	log.Debugf("Loaded configuration: %v", config)
	count, err := cmd.Flags().GetInt("count")
	if err != nil {
		log.Fatalf("could not parse --count flag: %v", err)
	}
	direction, err := cmd.Flags().GetString("direction")
	if err != nil {
		log.Fatalf("could not parse --direction flag: %v", err)
	}
	label := args[0]
	url := args[1]
	directory := label
	if len(args) == 3 {
		directory = args[2]
	}
	// check the label does not exist
	for _, podcast := range config.Podcasts {
		if podcast.Label == label {
			log.Errorf("a podcast with label %s already exists (for URL %s)", label, podcast.Feed)
			os.Exit(1)
		}
	}
	log.Infof("adding podcast: %s with feed %s to %d directory", label, url, directory)
	config.Podcasts = append(config.Podcasts, &feed.Podcast{
		Label:       label,
		Feed:        url,
		Directory:   directory,
		CountToKeep: count,
		Start:       direction,
		Episodes:    make(map[string]*feed.Episode, 0),
	})
	backend.Save(config)
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().IntP("count", "o", 0, "number of episodes to keep on disk, default is 0 which honors the master config default")
	addCmd.Flags().StringP("direction", "r", "oldest", "order of podcasts, 'oldest' or 'newest'")

}
