package migrate

import (
	"github.com/friendlycaptcha/friendly-stripe-sync/internal/entry/migrate"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var migratableStores = []string{"postgres"}

func Setup() *cobra.Command {
	migrateRootCmd := &cobra.Command{
		Use:   "migrate [store] [operation]",
		Short: "Migrate given data store",
	}

	for _, storeName := range migratableStores {
		c := buildMigrationCommand(storeName)
		c.PersistentFlags().Bool("danger", false, "Pass --danger to acknowledge a potentially dangerous operation.")
		migrateRootCmd.AddCommand(c)
	}

	return migrateRootCmd
}

func getVersionFlagValue(cmd *cobra.Command) int {
	v, err := cmd.Flags().GetInt("version")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get version flag")
	}
	return v
}

func buildMigrationCommand(datastoreName string) *cobra.Command {
	r := &cobra.Command{
		Use:   datastoreName + " [operation]",
		Short: "Migrate " + datastoreName + " with the given operation",
	}

	up := &cobra.Command{
		Use:   "up",
		Short: "Migrates the store to the latest version",
		Run: func(cmd *cobra.Command, args []string) {
			migrate.Migrate(datastoreName, "up", migrate.MigrateOpts{})
		},
	}

	down := &cobra.Command{
		Use:   "down",
		Short: "Migrates the store to the earliest version",
		Run: func(cmd *cobra.Command, args []string) {
			migrate.Migrate(datastoreName, "down", migrate.MigrateOpts{})
		},
	}
	down.Flags().Bool("danger", false, "Pass --danger to acknowledge this is potentially dangerous.")
	down.MarkFlagRequired("danger")

	version := &cobra.Command{
		Use:   "version",
		Short: "Prints the current version and \"dirty\" state",
		Run: func(cmd *cobra.Command, args []string) {
			migrate.Migrate(datastoreName, "version", migrate.MigrateOpts{})
		},
	}

	list := &cobra.Command{
		Use:   "list",
		Short: "Lists the migrations known to the application",
		Run: func(cmd *cobra.Command, args []string) {
			migrate.Migrate(datastoreName, "list", migrate.MigrateOpts{})
		},
	}

	force := &cobra.Command{
		Use:   "force",
		Short: "Forces the migration state to the given version",
		Run: func(cmd *cobra.Command, args []string) {
			migrate.Migrate(datastoreName, "force", migrate.MigrateOpts{
				TargetVersion: getVersionFlagValue(cmd),
			})
		},
	}
	force.Flags().Int("version", 9999, "Version to set the state to")
	force.MarkFlagRequired("version")
	force.Flags().Bool("danger", false, "Pass --danger to acknowledge this is potentially dangerous.")
	force.MarkFlagRequired("danger")

	to := &cobra.Command{
		Use:   "to",
		Short: "Migrates to the given version (up or down)",
		Run: func(cmd *cobra.Command, args []string) {
			migrate.Migrate(datastoreName, "to", migrate.MigrateOpts{
				TargetVersion: getVersionFlagValue(cmd),
			})
		},
	}
	to.Flags().Int("version", 9999, "Version to migrate to")
	to.MarkFlagRequired("version")

	to.MarkFlagRequired("danger")

	r.AddCommand(up)
	r.AddCommand(down)
	r.AddCommand(version)
	r.AddCommand(list)
	r.AddCommand(force)
	r.AddCommand(to)
	return r
}
