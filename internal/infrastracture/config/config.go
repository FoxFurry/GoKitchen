package config

import (
	"github.com/spf13/viper"
	"log"
)

func LoadConfig(){
	viper.AddConfigPath("./config")

	viper.SetConfigName("cfg")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("Could not read config file: %v", err)
	}
}
