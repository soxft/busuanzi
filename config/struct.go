package config

type Parser struct {
	Redis struct {
		Address   string `yaml:"Address"`
		Password  string `yaml:"Password"`
		Database  int    `yaml:"Database"`
		Prefix    string `yaml:"Prefix"`
		MaxIdle   int    `yaml:"MaxIdle"`
		MaxActive int    `yaml:"MaxActive"`
	} `yaml:"Redis"`
	Web struct {
		Address string `yaml:"Address"`
		AcAo    string `yaml:"Access-Control-Allow-Origin"`
		Debug   bool   `yaml:"Debug"`
		Log     bool   `yaml:"Log"`
	} `yaml:"Web"`
	Bsz struct {
		Expire int `yaml:"Expire"`
	} `yaml:"Bsz"`
}
