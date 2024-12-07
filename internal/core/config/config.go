package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"sync"
)

type Config struct {
	DiscordToken string
	DSN          string
}

var (
	once     sync.Once
	instance *Config
)

func Inst() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("No .env file found, loading from OS environment variables.")
		}

		instance = &Config{
			DiscordToken: getEnv("DISCORD_BOT_TOKEN", "discord"),
			DSN:          getEnv("DATABASE_URL", "dsn"),
		}
	})
	return instance
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
