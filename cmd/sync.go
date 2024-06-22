/*
Copyright © 2023 Daniel Blezek <blezek.daniel@mayo.edu>
This file is part of a CLI application.
*/
package cmd

import (
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Download and sync podcasts",
	Long: `Load the config file, fetch episodes from the RSS feed,
           compare to the files downloaded or deleted.  Updates files
           to keep the count of local files.`,
	Args: cobra.ExactArgs(0),
	Run:  Sync,
}

func Sync(cmd *cobra.Command, args []string) {
	backend, config := LoadConfiguration(cmd)
	for _, podcast := range config.Podcasts {
		podcast.Sync(config)
		log.Debugf("found podcast: %#v", spew.Sdump(podcast))
	}
	backend.Save(config)
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
