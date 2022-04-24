package config

type Parser struct {
	Redis struct {
		Address  string `yaml:"Address"`
		Password string `yaml:"Password"`
		Database int    `yaml:"Database"`
		Prefix   string `yaml:"Prefix"`
	}
	Web struct {
		Address string `yaml:"Address"`
		AcAo    string `yaml:"Access-Control-Allow-Origin"`
		Debug   bool   `yaml:"Debug"`
	}
	Bsz struct {
		Expire int `yaml:"Expire"`
	}
}
