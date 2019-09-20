package common

import (
	"fmt"

	"github.com/spf13/viper"
)

func ReadConfig(file string) error {
	if file != "" {
		viper.SetConfigFile(file)
	} else {
		viper.SetConfigName("app")
		viper.AddConfigPath(".")
	}

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config failed, %v", err)
	}

	return nil
}
