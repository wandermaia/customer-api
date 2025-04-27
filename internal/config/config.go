package config

import (
	"github.com/spf13/viper"
)

// Config contém todas as configurações da aplicação
type Config struct {
	ServerPort  string `mapstructure:"SERVER_PORT"`
	DBHost      string `mapstructure:"DB_HOST"`
	DBPort      string `mapstructure:"DB_PORT"`
	DBUser      string `mapstructure:"DB_USER"`
	DBPassword  string `mapstructure:"DB_PASSWORD"`
	DBName      string `mapstructure:"DB_NAME"`
	Environment string `mapstructure:"ENVIRONMENT"`
}

// LoadConfig carrega as configurações das variáveis de ambiente
func LoadConfig() (*Config, error) {

	viper.AutomaticEnv()

	config := &Config{
		ServerPort:  viper.GetString("SERVER_PORT"),
		DBHost:      viper.GetString("DB_HOST"),
		DBPort:      viper.GetString("DB_PORT"),
		DBUser:      viper.GetString("DB_USER"),
		DBPassword:  viper.GetString("DB_PASSWORD"),
		DBName:      viper.GetString("DB_NAME"),
		Environment: viper.GetString("ENVIRONMENT"),
	}

	// Valores padrão
	if config.ServerPort == "" {
		config.ServerPort = "8080"
	}
	if config.Environment == "" {
		config.Environment = "development"
	}

	return config, nil
}
