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
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	PersistentPreRun: PreRunRoot,
}

var Backend feed.BackendInterface

func PreRunRoot(cmd *cobra.Command, args []string) {
	debug, _ := cmd.Flags().GetBool("debug")
	if debug {
		log.SetLevel(log.DebugLevel)
	}
	//backend, _ := cmd.Flags().GetString("backend")
	//if backend == "db" {
	//	Backend = &feed.SqliteBackend{}
	//} else {
	//	Backend = &feed.FileBackend{}
	//}
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.castigate.yaml)")
	rootCmd.PersistentFlags().StringP("config", "c", "castigate.yaml", "path to config file")
	//rootCmd.PersistentFlags().StringP("backend", "e", "db", "backend configuration, should be 'db' or 'yaml'")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "set logging to debug")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
