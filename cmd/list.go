/*
Copyright © 2023 Daniel Blezek <blezek.daniel@mayo.edu>
This file is part of a CLI application.
*/
package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list podcasts",
	Long:  ``,
	Args:  cobra.NoArgs,
	Run:   RunListCmd,
}

func RunListCmd(cmd *cobra.Command, args []string) {
	_, config := LoadConfiguration(cmd)
	log.Debugf("Loaded configuration: %v", config)
	for _, podcast := range config.Podcasts {
		fmt.Printf("Label: %s\nFeed: %s\nDirection: %s\nNumber of Episodes: %d\n\n",
			podcast.Label, podcast.Feed, podcast.Start, len(podcast.Episodes))
	}
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
