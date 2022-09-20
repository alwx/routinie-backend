package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ReadConfig() {
	viper.SetDefault("environment", "development")
	viper.SetEnvPrefix("routinie")
	viper.BindEnv("environment")
	viper.AutomaticEnv()

	viper.SetConfigName(viper.GetString("environment"))
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error while loading config: %s", err)
	}
}
