package cfg

import "os"

type Config struct {
	Postgres PostgresConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func Load() Config {
	return Config{
		Postgres: PostgresConfig{
			Host:     getenv("POSTGRES_HOST", "localhost"),
			Port:     getenv("POSTGRES_PORT", "5432"),
			User:     getenv("POSTGRES_USER", "postgres"),
			Password: getenv("POSTGRES_PASSWORD", "postgres"),
			Database: getenv("POSTGRES_DB", "friday"),
		},
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
