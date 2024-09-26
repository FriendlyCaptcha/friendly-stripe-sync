package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/friendlycaptcha/friendly-stripe-sync/internal/config/cfgmodel"
)

var CfgFile string

func InitConfig() {
	if CfgFile != "" {
		// Use config file from the flag.
		fmt.Println("Using config file from flag:", CfgFile)
		viper.SetConfigFile(CfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".friendly" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".friendly-stripe-sync")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.SetEnvPrefix("fss")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("WARN could not find config file", err)
	} else {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	setupDefaults()
}

// GetStruct returns the config as a plain struct. You should call this after InitConfig.
func GetStruct() cfgmodel.FriendlyStripeSync {
	cfg := cfgmodel.FriendlyStripeSync{}

	err := viper.Unmarshal(&cfg)
	if err != nil {
		panic("Failed to unmarshal config: " + err.Error())
	}

	return cfg
}
