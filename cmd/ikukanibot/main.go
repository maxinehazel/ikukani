package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"git.maxinekrebs.dev/softpunk/ikukani"
	"git.maxinekrebs.dev/softpunk/ikukani/internal/notifier"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile          string
	waniKaniToken    string
	twilioFromNumber string
	twilioToNumber   string
	twilioToken      string
	twilioAccountSID string

	rootCmd = &cobra.Command{
		Use:   "ikukanibot",
		Short: "cli for interfacing with the WaniKani API",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ikukani.Token = viper.GetString("wk_token")
		},
	}

	reviewCmd = &cobra.Command{
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
)

var nextReviewCmd = &cobra.Command{
	Use:   "next",
	Short: "Print when the next review is availabe",
	Run: func(cmd *cobra.Command, args []string) {
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
		resp, err := ikukani.ReviewAvailable()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp)
	},
}

var notifyReviewCmd = &cobra.Command{
	Use:   "notify",
	Short: "sends text notification when review is ready",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ikukani.Token = viper.GetString("wk_token")
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("listening for available reviews")
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs)
		go func() {
			s := <-sigs
			log.Printf("RECEIVED SIGNAL: %s", s)
			os.Exit(1)
		}()
		var waitInterval time.Duration
		var sent bool
		var onVacation bool
		var err error
		for {
			onVacation, err = ikukani.VacationMode()
			if err != nil {
				log.Fatal(err)
			}

			if onVacation == true {
				log.Println("User on vacation, checking again later")
				waitInterval = time.Hour * 24
			} else {
				resp, err := ikukani.ReviewAvailable()
				if err != nil {
					log.Fatal(err)
				}

				if resp == true && sent == false {
					log.Println("review available, sending text notification")
					notifier.TwilioSID = viper.GetString("twilio_sid")
					notifier.TwilioToken = viper.GetString("twilio_token")
					n := notifier.Notification{
						From: viper.GetString("twilio_from"),
						To:   viper.GetString("twilio_to"),
						Body: "You have reviews ready! https://www.wanikani.com/dashboard",
					}
					sid, err := n.Send()
					if err != nil {
						log.Fatal(err)
					}
					log.Println("Message sid: " + sid)
					sent = true
					waitInterval = time.Minute * 30
				} else if resp == true && sent == true {
					log.Println("review available but incomplete, checking again later")
				} else {
					log.Println("review not available, checking again later")
					nextReview, err := ikukani.NextReviewInDuration()
					if err != nil {
						log.Fatal(err)
					}
					waitInterval = *nextReview
					sent = false
				}
			}

			log.Println("waiting for " + waitInterval.String())
			time.Sleep(waitInterval)
		}
	},
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ikukani.yaml)")
	rootCmd.PersistentFlags().StringVarP(&waniKaniToken, "token", "t", "", "WaniKani API v2 token (required)")
	notifyReviewCmd.Flags().StringVarP(&twilioAccountSID, "sid", "s", "", "Twilio account SID")
	notifyReviewCmd.Flags().StringVarP(&twilioToken, "twtoken", "w", "", "Twilio account Token")
	notifyReviewCmd.Flags().StringVarP(&twilioFromNumber, "from", "f", "", "Number to send notifications from")
	notifyReviewCmd.Flags().StringVarP(&twilioToNumber, "to", "T", "", "Number to send notifications to")

	viper.BindPFlag("wk_token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("twilio_sid", notifyReviewCmd.Flags().Lookup("sid"))
	viper.BindPFlag("twilio_token", notifyReviewCmd.Flags().Lookup("twtoken"))
	viper.BindPFlag("twilio_from", notifyReviewCmd.Flags().Lookup("from"))
	viper.BindPFlag("twilio_to", notifyReviewCmd.Flags().Lookup("to"))

	rootCmd.AddCommand(reviewCmd)
	reviewCmd.AddCommand(nextReviewCmd, availabeReviewCmd, notifyReviewCmd)
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
