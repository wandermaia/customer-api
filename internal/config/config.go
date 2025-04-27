package config

import (
	"fmt"

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
	// Tenta carregar o arquivo .env na raiz do projeto
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Tenta ler o arquivo .env, mas não retorna erro se não encontrar
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Arquivo .env não encontrado, isso é aceitável
			fmt.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
		} else {
			// Outro erro ao ler o arquivo
			fmt.Printf("Erro ao ler arquivo .env: %s\n", err)
		}
	} else {
		fmt.Println("Usando configurações do arquivo .env")
	}

	// Permite que as variáveis de ambiente do sistema substituam as do arquivo
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
