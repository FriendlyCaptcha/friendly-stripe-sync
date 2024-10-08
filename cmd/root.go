package cmd

import (
	"context"
	"io"

	"github.com/friendlycaptcha/friendly-stripe-sync/cmd/migrate"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/buildinfo"
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:              "friendly-stripe-sync",
	Short:            "Load data and synchronize events from Stripe to Postgres",
	Long:             `Load data and synchronize events from Stripe to Postgres`,
	PersistentPreRun: bindFlags,
}

func Execute(ctx context.Context, stdout io.Writer, args []string) error {
	rootCmd.SetArgs(args)
	rootCmd.SetOut(stdout)
	rootCmd.SilenceUsage = true

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.PersistentFlags().StringVar(&config.CfgFile, "config", "", "config file (default is $HOME/.friendly-stripe-sync.yaml)")

	rootCmd.Version = buildinfo.FullVersion()

	rootCmd.PersistentFlags().BoolP("development", "d", false, "Development mode (prints prettier log messages)")
	rootCmd.PersistentFlags().BoolP("debug", "D", false, "Debug mode (prints debug messages and call traces)")
	rootCmd.AddCommand(migrate.Setup())
}

func bindFlags(cmd *cobra.Command, args []string) {
	viper.BindPFlag("development", cmd.Flags().Lookup("development"))
	viper.BindPFlag("debug", cmd.Flags().Lookup("debug"))
}
