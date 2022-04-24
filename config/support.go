package config

type ConfigStruct struct {
	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		Database int    `yaml:"database"`
		Prefix   string `yaml:"prefix"`
	}
	Web struct {
		Address string `yaml:"address"`
	}
}
