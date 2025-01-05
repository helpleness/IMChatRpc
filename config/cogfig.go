package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func ConfigInit() {
	workDir, _ := os.Getwd()

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
