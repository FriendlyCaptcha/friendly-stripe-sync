package cmd

import (
	"github.com/friendlycaptcha/friendly-stripe-sync/entry/load"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Loads all data from Stripe and puts it in the DB",
	Long:  "Loads all data from Stripe and puts it in the DB",
	RunE: func(cmd *cobra.Command, _ []string) error {
		cfg := config.GetStruct()
		return load.Start(cmd.Context(), cfg)
	},
	PreRun: bindLoad,
}

func init() {
	rootCmd.AddCommand(loadCmd)

	loadCmd.Flags().Bool("purge", false, "Delete all existing data before loading")
}

func bindLoad(cmd *cobra.Command, args []string) {
	config.InitConfig()
	viper.BindPFlag("purge", cmd.Flags().Lookup("purge"))
}
