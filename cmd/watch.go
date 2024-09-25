package cmd

import (
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/config"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/entry/watch"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Sync events from Stripe into the database on a fixed interval",
	Long:  "Load the initial data if it hasn't been loaded already and periodically load events from Stripe to the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.GetStruct()
		return watch.Start(cmd.Context(), cfg)
	},
	PreRun: bindWatch,
}

func init() {
	rootCmd.AddCommand(watchCmd)

	watchCmd.Flags().IntP("interval_seconds", "i", 60, "The seconds to wait between each sync")
}

func bindWatch(cmd *cobra.Command, args []string) {
	config.InitConfig()

	viper.BindPFlag("stripe_sync.interval_seconds", cmd.Flags().Lookup("interval_seconds"))
}
