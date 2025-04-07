package config

import "os"

type Config struct {
	Addr             string
	OrdersCollection string
	DBname           string
	Uri              string
	JWTSecret        string
}

func NewConfig() *Config {
	return &Config{
		Addr:             getEnv("ADDR", "8080"),
		Uri:              getEnv("MONGO", ""),
		DBname:           getEnv("DBNAME", "orders"),
		OrdersCollection: getEnv("ORDERS_COLLECTION", "orders"),
		JWTSecret:        getEnv("JWT_SECRET", "secret"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
