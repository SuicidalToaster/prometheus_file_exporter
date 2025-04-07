package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

type Config struct {
	WatchedDirs  []string
	WatchedFiles []string
}

func InitConfig() *viper.Viper {
	var err error
	cfg := viper.New()
	dirname, err := os.UserHomeDir()
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatal(err)
		}
		cfg.AddConfigPath(dirname)
	}
	cfg.SetEnvPrefix("PFE_")
	cfg.AutomaticEnv()
	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.AddConfigPath(".")
	cfg.AddConfigPath("~")
	cfg.AddConfigPath("$HOME")
	cfg.AddConfigPath("/etc/file_exporter")
	cfg.AddConfigPath("$HOME/.config")
	cfg.SetConfigName("pfe.config")
	cfg.SetConfigType("yaml")
	cfg.Get("hash")
	cfg.SetDefault("DirPaths", []string{})
	cfg.SetDefault("HashFiles", []string{})
	err = cfg.ReadInConfig()

	log.Print(cfg.ConfigFileUsed())
	log.Print(cfg.AllSettings())
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		} else {
			panic(err.Error())
		}
	}
	return cfg
}
