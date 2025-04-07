package config

import "os"

type Config struct {
	Addr string
	ComicsCollection string
	DBname string
	Uri string
}

func NewConfig() *Config {
	return &Config{
		Addr: getEnv("ADDR", ":8080"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
