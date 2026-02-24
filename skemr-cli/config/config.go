package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initConfig() {
	viper.SetDefault("controlPlaneUrl", "https://api.skemr.com")

	viper.SetEnvPrefix("SKEMR")
	viper.AutomaticEnv()
}

func init() {
	cobra.OnInitialize(initConfig)
}
