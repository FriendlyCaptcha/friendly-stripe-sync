package cmd

import (
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/config"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/entry/sync"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v74"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Loads changes from the last 30 days from Stripe and applies them to the DB",
	Long:  "Loads changes from the last 30 days from Stripe and applies them to the DB",
	Run: func(cmd *cobra.Command, args []string) {
		sync.Start()
	},
	PreRun: bindSync,
}

func init() {
	rootCmd.AddCommand(syncCmd)
}

func bindSync(cmd *cobra.Command, args []string) {
	config.InitConfig()

	stripe.Key = viper.GetString("stripe.api_key")
}
