package config

import (
	"flag"
	"log"
	"strings"

	"github.com/spf13/viper"
)

var (
	configPath string
	DistPath   string
	VERSION    = "2.8.8"
	DEBUG      bool
)

func Init() {
	// get config file path
	flag.StringVar(&configPath, "c", "config.yaml", "config path")
	flag.StringVar(&DistPath, "d", "dist", "dist path")
	flag.Parse()

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("read config file error: %s", err)
	}

	log.Printf("[INFO] Config loaded %s", viper.AllSettings())

	DEBUG = viper.GetBool("debug")
}
