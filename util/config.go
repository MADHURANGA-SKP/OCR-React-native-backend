package util

import (
	"github.com/spf13/viper"
)

// config stores all the configuration details of the program
// the values are read by viper  from config file or enviroment vaiable
type Config struct {
	Environment       string `mapstructure:"ENVIRONMENT"`
	DBSource          string `mapstructure:"DB_SOURCE"`
	MigrationURL      string `mapstructure:"MIGRATION_URL"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	FrontEndOrigin    string `mapstructure:"FRONT_END_ORIGIN"`
}

// loadconfig reads configurations from file or enviornment variables
func Loadconfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
