package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Load loads configuration values.
// It will try to load a .env file (if present) into environment variables for local development,
// then it will initialize Viper to read from environment variables and optional config files.
func Load() {
	// Load .env if present (don't fail if not present)
	_ = godotenv.Load()

	viper.AutomaticEnv()

	// Allow optional config file (yaml/json/toml) named "config" in current folder or ./config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Loaded config file: %s\n", viper.ConfigFileUsed())
	}
}

// Helper getters
func GetString(key string) string {
	return viper.GetString(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}
