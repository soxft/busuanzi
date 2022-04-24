package config

import (
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	C *ConfigStruct
)

func Init() {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic("Error reading config file")
	}
	C = &ConfigStruct{}
	err = yaml.Unmarshal(data, C)
	if err != nil {
		panic("Error parsing config file")
	}
}
