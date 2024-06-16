/*
Copyright © 2023 Daniel Blezek <blezek.daniel@mayo.edu>
This file is part of a CLI application.
*/
package cmd

import (
	"castigate/castigate"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: Sync,
}

func Sync(cmd *cobra.Command, args []string) {
	configFile := args[0]
	log.Infof("loading %s", configFile)
	// read the yaml file
	b, err := os.ReadFile(args[0])
	if err != nil {
		log.Fatalf("failed to read input: %s", err)
	}
	config := castigate.Config{
		Podcasts:           make([]*castigate.Podcast, 0),
		FilenameTemplate:   `{{.episode.Date.Format "2006-01-02-15:04:05" }}-{{.item.Title}}.mp3`,
		DefaultCountToKeep: 10,
	}

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		log.Fatalf("failed to load YAML: %s", err)
	}
	for _, podcast := range config.Podcasts {
		podcast.Sync(config)
		log.Debugf("found podcast: %#v", spew.Sdump(podcast))
	}
	log.Infof("saving backup to %s", configFile+".bak")
	err = os.WriteFile(configFile+".bak", buffer, 0644)
	if err != nil {
		log.Fatalf("failed to write backup config: %s", err)
	}

	buffer, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalf("failed to marshal YAML: %s", err)
	}

	// Write the new config
	err = os.WriteFile(configFile, buffer, 0644)
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
