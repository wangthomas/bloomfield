package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

var Config Configuration

type Configuration struct {
	Port       string 
}



func LoadConfig(configFile string) {

	if _, err := toml.DecodeFile(configFile, &Config); err != nil {
		fmt.Println("Cannot load config file. Using default config.", err)
		LoadDefault()
		return
	}
}

func LoadDefault() {
	Config.Port = "8679"
}

