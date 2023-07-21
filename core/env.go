package core

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Env()  {
	godotenv.Load()
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	settings()
}

func settings()  {
	viper.SetDefault("APP_URL", "localhost")
	viper.SetDefault("APP_PORT", 80)
	viper.SetDefault("DEBUG", true)
	viper.SetDefault("CONCURRENCY_LIMIT", 1024)
	viper.SetDefault("REQUEST_TIMEOUT", 30)
	viper.SetDefault("REQUEST_SIZE_LIMIT", 1024)

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 3306)
	viper.SetDefault("DB_USER", "root")
	viper.SetDefault("DB_PASS", "123456")
	viper.SetDefault("DB_NAME", "db")

	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PASS", "123456")
	viper.SetDefault("REDIS_PORT", 6379)
	viper.SetDefault("REDIS_DB", 0)
}
