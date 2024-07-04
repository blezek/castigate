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
		fmt.Print(podcast.PrintDetails())
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
}
