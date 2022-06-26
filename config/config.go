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
	DistPath   string
)

func init() {
	// get config file path
	flag.StringVar(&configPath, "c", "config.yaml", "config path")
	flag.StringVar(&DistPath, "d", "dist", "dist path")
	flag.Parse()

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal("Error reading config file:\r\n" + err.Error())
	}
	C = &Config{}
	err = yaml.Unmarshal(data, C)
	if err != nil {
		log.Fatal("Error parsing config file:\r\n" + err.Error())
	}

	Redis = C.Redis
	Web = C.Web
	Bsz = C.Bsz
}
