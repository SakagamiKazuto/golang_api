package config

import "github.com/BurntSushi/toml"


type DBConfig struct {
	Name     string `toml:"name"`
	Host     string `toml:"host"`
	User     string `toml:"user"`
	Port     int    `toml:"port"`
	Password string `toml:"pass"`
}

type Config struct {
	DB DBConfig `toml:"database"`
}

func NewConfig(path string, appMode string) (Config, error) {
	var conf Config

	confPath := path + appMode + ".toml"
	if _, err := toml.DecodeFile(confPath, &conf); err != nil {
		return conf, err
	}

	return conf, nil
}
