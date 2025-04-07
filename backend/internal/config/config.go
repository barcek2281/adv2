package config

import "os"

type Config struct {
	Addr             string
	ComicsCollection string
	DBname           string
	Uri              string
	JWTSecret string
}

func NewConfig() *Config {
	return &Config{
		Addr: getEnv("ADDR", ":8080"),
		Uri: getEnv("MONGO", ""),
		DBname: getEnv("DBname", "product"),
		ComicsCollection: getEnv("COMICS_COLLECTION", "comics"),
		JWTSecret: getEnv("JWTSecret", "secret"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
