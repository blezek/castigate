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
	backend.Save(config)
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
