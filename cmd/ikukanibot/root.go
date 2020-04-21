package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/softpunks/ikukani"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

const (
	apiVersion = "20170710"
)

var (
	// Used for flags.
	cfgFile       string
	waniKaniToken string

	// Client for making requests to wk api
	wkClient      *ikukani.Client

	rootCmd = &cobra.Command{
		Use:   "ikukanibot",
		Short: "cli for interfacing with the WaniKani API",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			apiToken := viper.GetString("wk_token")
			wkClient = ikukani.NewClient(apiToken, apiVersion)
		},
	}
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ikukani.yaml)")
	rootCmd.PersistentFlags().StringVarP(&waniKaniToken, "token", "t", "", "WaniKani API v2 token (required)")

	err := viper.BindPFlag("wk_token", rootCmd.PersistentFlags().Lookup("token"))

	if err != nil {
		log.Fatal(err)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".ikukani")
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err == nil {
		// fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
