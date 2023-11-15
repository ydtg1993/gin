package core

import (
	"fmt"
	"github.com/spf13/viper"
)

var Config *viper.Viper

func init() {
	Config = viper.New()
	Config.SetConfigName("config")
	Config.AddConfigPath(".")
	Config.SetConfigType("yaml")

	if err := Config.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("Error reading config file: %s\n", err))
	}
}
