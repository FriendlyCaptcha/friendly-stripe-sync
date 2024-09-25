package cmd

import (
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/config"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/entry/sync"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Loads changes from the last 30 days from Stripe and applies them to the DB",
	Long:  "Loads changes from the last 30 days from Stripe and applies them to the DB",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.GetStruct()
		return sync.Start(cmd.Context(), cfg)
	},
	PreRun: bindSync,
}

func init() {
	rootCmd.AddCommand(syncCmd)
}

func bindSync(cmd *cobra.Command, args []string) {
	config.InitConfig()
}
