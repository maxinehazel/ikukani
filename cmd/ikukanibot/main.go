package main

import (
	"fmt"
	"log"
	"os"

	"git.maxinekrebs.dev/softpunk/ikukani"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile       string
	waniKaniToken string

	rootCmd = &cobra.Command{
		Use:   "ikukanibot",
		Short: "cli for interfacing with the WaniKani API",
	}

	reviewCmd = &cobra.Command{
		Use:   "review",
		Short: "interact with review data",
	}
)

var nextReviewCmd = &cobra.Command{
	Use:   "next",
	Short: "Print when the next review is availabe",
	Run: func(cmd *cobra.Command, args []string) {
		ikukani.Token = viper.GetString("wk_token")
		resp, err := ikukani.NextReviewIn()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp)
	},
}

var availabeReviewCmd = &cobra.Command{
	Use:   "available",
	Short: "Returns is there are any available reviews",
	Run: func(cmd *cobra.Command, args []string) {
		ikukani.Token = viper.GetString("wk_token")
		resp, err := ikukani.ReviewAvailable()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp)
		if resp == true {
			// send text message
		}
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ikukani.yaml)")
	rootCmd.PersistentFlags().StringVarP(&waniKaniToken, "token", "t", "", "WaniKani API v2 token (required)")

	viper.BindPFlag("wk_token", rootCmd.PersistentFlags().Lookup("token"))

	rootCmd.AddCommand(reviewCmd)
	reviewCmd.AddCommand(nextReviewCmd, availabeReviewCmd)
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

func main() {
	rootCmd.Execute()
}
