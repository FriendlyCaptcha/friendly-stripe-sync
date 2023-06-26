package cmd

import (
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/config"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/entry/load"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v74"
)

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Loads all data from Stripe and puts it in the DB",
	Long:  "Loads all data from Stripe and puts it in the DB",
	Run: func(cmd *cobra.Command, args []string) {
		load.Start()
	},
	PreRun: bindLoad,
}

func init() {
	rootCmd.AddCommand(loadCmd)

	loadCmd.Flags().Bool("purge", false, "Delete all existing data before loading")
}

func bindLoad(cmd *cobra.Command, args []string) {
	config.InitConfig()
	stripe.Key = viper.GetString("stripe.api_key")

	viper.BindPFlag("purge", cmd.Flags().Lookup("purge"))
}
