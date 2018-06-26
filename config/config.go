package config

import (
	"github.com/spf13/viper"
	"fmt"
)

const configPath = "config/app.yaml"

var C *viper.Viper

func init() {
	C = NewConfig()
}

func NewConfig() *viper.Viper {
	if C != nil {
		return C
	}

	C = viper.New()
	C.SetConfigType("yaml")
	C.SetConfigFile(configPath)
	err := C.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	return C
}