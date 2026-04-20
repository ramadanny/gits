package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func Init() {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".config", "ramadanny-gits")
	os.MkdirAll(configPath, os.ModePerm)

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	if err := viper.ReadInConfig(); err != nil {
		viper.WriteConfigAs(filepath.Join(configPath, "config.json"))
	}
}

func Set(key, value string) {
	viper.Set(key, value)
	viper.WriteConfig()
}

func Get(key string) string {
	return viper.GetString(key)
}