package config

import (
	"medbuddy-backend/utility"
	"os"

	"github.com/spf13/viper"
)

type Configuration struct {
	ServerPort string `mapstructure:"SERVER_PORT"`
	SecretKey  string `mapstructure:"SECRET_KEY"`
	MongoHost  string `mapstructure:"MONGO_HOST"`
}

// Setup initialize configuration
var (
	Config *Configuration
)

func Setup() {
	var configuration *Configuration
	logger := utility.NewLogger()

	viper.SetConfigName("sample")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Overwrite file env's from environment
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logger.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		logger.Fatalf("Unable to decode into struct, %v", err)
	}

	if port := os.Getenv("PORT"); port != "" {
		configuration.ServerPort = port
	}

	Config = configuration
	logger.Info("CONFIGURATIONS LOADED SUCCESSFULLY")
}

// GetConfig helps you to get configuration data
func GetConfig() *Configuration {
	return Config
}
