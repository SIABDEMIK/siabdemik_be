package config

import (
	"Ayala-Crea/server-app-absensi/models"
	"fmt"

	"github.com/spf13/viper"
)

func LoadConfig() (*models.Config, error) {
	var Config models.Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &Config, nil
}
