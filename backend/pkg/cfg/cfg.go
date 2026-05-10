package cfg

import "os"

type Config struct {
	HTTPAddr string
	Postgres PostgresConfig
	Resend   ResendConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type ResendConfig struct {
	APIKey string
	From   string
}

func Load() Config {
	return Config{
		HTTPAddr: getenv("HTTP_ADDR", ":8080"),
		Postgres: PostgresConfig{
			Host:     getenv("POSTGRES_HOST", "localhost"),
			Port:     getenv("POSTGRES_PORT", "5432"),
			User:     getenv("POSTGRES_USER", "postgres"),
			Password: getenv("POSTGRES_PASSWORD", "postgres"),
			Database: getenv("POSTGRES_DB", "friday"),
		},
		Resend: ResendConfig{
			APIKey: getenv("RESEND_API_KEY", ""),
			From:   getenv("RESEND_FROM", "onboarding@resend.dev"),
		},
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
