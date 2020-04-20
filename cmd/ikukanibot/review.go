package main

import (
	"fmt"
	"log"
	"os"

	"github.com/softpunks/ikukani"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "interact with review data",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ikukani.Token = viper.GetString("wk_token")
		vacation, err := ikukani.VacationMode()
		if err != nil {
			log.Fatal(err)
		}

		if vacation {
			fmt.Println("User on vacation")
			os.Exit(0)
		}
	},
}

var nextReviewCmd = &cobra.Command{
	Use:   "next",
	Short: "Print when the next review is availabe",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := ikukani.NextReviewInString()
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
		resp, err := ikukani.ReviewAvailable()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp)
	},
}

func init() {
	rootCmd.AddCommand(reviewCmd)
	reviewCmd.AddCommand(nextReviewCmd, availabeReviewCmd)
}
