/*
Copyright © 2023 Daniel Blezek <blezek.daniel@mayo.edu>
This file is part of a CLI application.
*/
package cmd

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit a podcast",
	Long: `The edit command supports changing the feed URL, directory, start direction
the count to keep, and resetting all the episodes to a given state.`,
	Args: cobra.ExactArgs(1),
	Run:  runEditCmd,
}

func runEditCmd(cmd *cobra.Command, args []string) {
	backend, config := LoadConfiguration(cmd)

	label := args[0]
	podcast, err := config.FindPodcast(label)
	if err != nil {
		log.Fatalf("could not find podcast with label %s: %v", label, err)
	}
	log.Infof("editing label %s: %s", label, podcast.Title)
	url, err := cmd.Flags().GetString("url")
	if err != nil {
		log.Fatalf("could not get url flag %v", err)
	}
	if url != "" {
		podcast.Feed = url
	}

	count, err := cmd.Flags().GetInt("count")
	if err != nil {
		log.Fatalf("could not get count flag %v", err)
	}
	if count >= 0 {
		podcast.CountToKeep = count
	}

	directory, err := cmd.Flags().GetString("directory")
	if err != nil {
		log.Fatalf("could not get directory flag %v", err)
	}
	if directory != "" {
		podcast.Directory = directory
	}

	start, err := cmd.Flags().GetString("start")
	if err != nil {
		log.Fatalf("could not get start flag %v", err)
	}
	if start != "" {
		if start == "oldest" {
			podcast.Start = "oldest"
		} else {
			podcast.Start = "newest"
		}
	}
	log.Infof("saving configuration")
	backend.Save(config)
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().String("url", "", "URL of the podcast")
	editCmd.Flags().Int("count", -1, "Number of episodes to keep on disk")
	editCmd.Flags().String("directory", "", "Directory of the podcast")
	editCmd.Flags().String("start", "", "download starting with oldest or newest")
}
