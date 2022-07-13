package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var (
	C     *Config
	Redis RedisConfig
	Web   WebConfig
	Bsz   BszConfig
)

var (
	configPath string
	AssetsPath string
)

func init() {
	// get config file path
	flag.StringVar(&configPath, "c", "config.yaml", "config path")
	flag.StringVar(&AssetsPath, "d", "assets", "assets path")
	flag.Parse()

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("Error reading config file:\r\n" + err.Error())
	}
	C = &Config{}
	if err = yaml.Unmarshal(data, C); err != nil {
		log.Fatal("Error parsing config file:\r\n" + err.Error())
	}

	Redis = C.Redis
	Web = C.Web
	Bsz = C.Bsz
}
