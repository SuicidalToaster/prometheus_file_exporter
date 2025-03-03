package config

import (
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	WatchedDirs  []string
	WatchedFiles []string
}

func InitConfig() *viper.Viper {
	cfg := viper.New()
	cfg.SetEnvPrefix("PFE_")
	cfg.AutomaticEnv()
	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.AddConfigPath(".")
	cfg.AddConfigPath("~")
	cfg.AddConfigPath("~/.config")
	cfg.AddConfigPath("$HOME")
	cfg.AddConfigPath("/etc/file_exporter")
	cfg.AddConfigPath("$HOME/.config")
	cfg.SetConfigName("file_exporter")
	cfg.SetConfigType("yaml")
	cfg.SetDefault("DirPaths", []string{"$HOME"})
	err := cfg.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {

		} else {
			panic(err.Error())
		}
	}
	return cfg
}
