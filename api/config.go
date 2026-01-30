package api

import "github.com/spf13/viper"

type Config struct {
	Port         int
	Cors         bool
	AllowedHosts []string
}

func InitConfig() (*Config, error) {
	config := &Config{
		Port:         viper.GetInt("api.Port"),
		Cors:         viper.GetBool("api.Cors"),
		AllowedHosts: viper.GetStringSlice("api.AllowedHosts"),
	}
	if config.Port == 0 {
		config.Port = 8081
	}
	if len(config.AllowedHosts) == 0 {
		config.AllowedHosts = append(config.AllowedHosts, "*")
	}
	return config, nil
}
