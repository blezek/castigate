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

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize the config file",
	Long:  ``,
	Run:   runInitCmd,
}

func runInitCmd(cmd *cobra.Command, args []string) {
	log.Debug("init called")
	filenameTemplate, err := cmd.Flags().GetString("template")
	if err != nil {
		log.Fatalf("error reading template flag: %v", err)
	}
	count, err := cmd.Flags().GetInt("count")
	if err != nil {
		log.Fatalf("error reading count flag: %v", err)
	}
	config := feed.NewConfig()
	config.FilenameTemplate = filenameTemplate
	config.DefaultCountToKeep = count

	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		log.Fatalf("could not get config file: %v", err)
	}
	backend := feed.FileBackend{}
	backend.Init(configFile)
	backend.Save(config)
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("template", "t", feed.DefaultFilenameTemplate, "template for filenames")
	initCmd.Flags().Int("count", 10, "number of episodes to keep by default")
}
