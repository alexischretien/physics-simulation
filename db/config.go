package db

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DbUser     string
	DbPassword string
	DbAddress  string
	DbName     string
}

func InitConfig() (*Config, error) {
	config := &Config{
		DbUser:     viper.GetString("db.user"),
		DbPassword: viper.GetString("db.password"),
		DbAddress:  viper.GetString("db.address"),
		DbName:     viper.GetString("db.schema.name"),
	}
	if config.DbUser == "" {
		return nil, fmt.Errorf("Database user name must be set")
	}
	if config.DbAddress == "" {
		return nil, fmt.Errorf("Database Address must be set")
	}
	if config.DbName == "" {
		return nil, fmt.Errorf("Database schema name must be set")
	}
	return config, nil
}
