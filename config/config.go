package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var (
	C *Parser
)

func init() {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("Error reading config file:\r\n" + err.Error())
	}
	C = &Parser{}
	err = yaml.Unmarshal(data, C)
	if err != nil {
		log.Fatal("Error parsing config file:\r\n" + err.Error())
	}
}
