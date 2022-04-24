package config

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	C *Parser
)

func Init() {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic("Error reading config file")
	}
	err = yaml.Unmarshal(data, &Parser{})
	if err != nil {
		panic("Error parsing config file")
	}
}
