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
	Use:   "reviews",
	Short: "interact with review data",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		apiToken := viper.GetString("wk_token")
		wkClient = ikukani.NewClient(apiToken, apiVersion)

		vacation, err := wkClient.VacationMode()
		if err != nil {
			log.Fatal(err)
		}

		if vacation {
			fmt.Println("User on vacation")
			os.Exit(0)
		}
	},
}

var getReviewsCmd = &cobra.Command{
	Use:   "list",
	Short: "Print all reviews",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := wkClient.GetReviews()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp)
	},
}

var nextReviewCmd = &cobra.Command{
	Use:   "in",
	Short: "Print when the next review is available",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := wkClient.NextReviewInString()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("review available in", resp)
	},
}

var availableReviewCmd = &cobra.Command{
	Use:   "available",
	Short: "Returns is there are any available reviews",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := wkClient.ReviewAvailable()
		if err != nil {
			log.Fatal(err)
		}
		if resp {
			fmt.Println("reviews are available!")
		} else {
			fmt.Println("reviews not available yet :(")
		}
	},
}

func init() {
	rootCmd.AddCommand(reviewCmd)
	reviewCmd.AddCommand(nextReviewCmd, availableReviewCmd, getReviewsCmd)
}
