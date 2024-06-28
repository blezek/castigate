/*
Copyright © 2023 Daniel Blezek <blezek.daniel@mayo.edu>
This file is part of a CLI application.
*/
package cmd

import (
	"castigate/feed"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "castigate",
	Short: "A simple podCAST GATEway application",
	Long: `castigate or the podCAST GATEway, is a Go project to manage Podcasts for
portable MP3 players, or simply to download episodes of
podcasts from the command line.  Each podcast has a label and the following settings:

feed:         the URL of the RSS feed for the podcast
directory:    the directory where the podcast is stored, relative to the config file
counttokeep:  number of episodes to keep on disk, if 0, use the default configuration
start:        download the oldest or newest podcasts first, in order

A podcast has a number of episodes.  Each episode goes through three states, new, downloaded, and deleted.
When an episode appears on the feed, it is added to the podcast in a new state.  If the number of episodes
in the directory is less than counttokeep, episodes in the new state are downloaded and their state is 
updated to downloaded.  Finally, when a downloaded episode no longer exists on disk, it is moved to
the deleted state.

To manage podcasts, run 'castigate sync' to download new episodes from each podcast, and copy to
an MP3 player (or wherever).  When an episode has been played, simply delete it from the directory
and run 'castigate sync' to retrieve any new episodes.
`,

	PersistentPreRun: PreRunRoot,
}

var Backend feed.BackendInterface

func PreRunRoot(cmd *cobra.Command, args []string) {
	debug, _ := cmd.Flags().GetBool("debug")
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringP("config", "c", "castigate.yaml", "path to config file")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "set logging to debug")
}

func LoadConfiguration(cmd *cobra.Command) (feed.FileBackend, feed.Config) {
	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		log.Fatalf("could not get config file: %v", err)
	}
	var backend feed.FileBackend
	backend.Init(configFile)
	config, err := backend.Load()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}
	return backend, config
}
