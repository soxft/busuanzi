package config

import (
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

func init() {
	data, err := ioutil.ReadFile("config.yaml")
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
