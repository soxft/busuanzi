package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	C *Parser
)

func init() {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic("Error reading config file:\r\n" + err.Error())
	}
	C = &Parser{}
	err = yaml.Unmarshal(data, C)
	if err != nil {
		panic("Error parsing config file:\r\n" + err.Error())
	}
}
