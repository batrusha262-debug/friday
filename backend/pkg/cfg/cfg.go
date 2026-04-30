package cfg

import "os"

type Config struct {
	HTTPAddr string
	Postgres PostgresConfig
	SMTP     SMTPConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
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
		SMTP: SMTPConfig{
			Host:     getenv("SMTP_HOST", "smtp.gmail.com"),
			Port:     getenv("SMTP_PORT", "587"),
			Username: getenv("SMTP_USERNAME", ""),
			Password: getenv("SMTP_PASSWORD", ""),
			From:     getenv("SMTP_FROM", "noreply@friday"),
		},
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
