package util

import (
	"fmt"

	"github.com/spf13/viper"
)

func Loadconfig(path string) error {
	viper.SetConfigName("test") // name of config file (without extension)
	viper.SetConfigType("env")  // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(path)   // path to look for the config file in

	viper.AutomaticEnv()

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		return (fmt.Errorf("fatal error config file: %w", err))
	}
	return nil
}
