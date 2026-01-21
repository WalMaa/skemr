package config

import (
	"log/slog"

	"github.com/spf13/viper"
)

type AppConfig struct {
	ControlPlaneUrl string
}

var Cfg AppConfig

func init() {
	viper.SetDefault("controlPlaneUrl", "http://localhost:8080")
	viper.BindEnv("controlPlaneUrl", "URL")
	if err := viper.Unmarshal(&Cfg); err != nil {
		slog.Error("Unable to unmarshal config", err)
	}

}
