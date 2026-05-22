package cfg

import "os"

type Config struct {
	HTTPAddr string
	Postgres PostgresConfig
	Resend   ResendConfig
	SMTP     SMTPConfig
}

type PostgresConfig struct {
	DSN      string
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

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func Load() Config {
	return Config{
		HTTPAddr: httpAddr(),
		Postgres: PostgresConfig{
			DSN:      firstEnv("DATABASE_URL", "POSTGRES_URL"),
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
		SMTP: SMTPConfig{
			Host:     getenv("SMTP_HOST", ""),
			Port:     getenv("SMTP_PORT", "587"),
			Username: getenv("SMTP_USERNAME", ""),
			Password: getenv("SMTP_PASSWORD", ""),
			From:     getenv("SMTP_FROM", getenv("SMTP_USERNAME", "")),
		},
	}
}

func httpAddr() string {
	if v := os.Getenv("HTTP_ADDR"); v != "" {
		return v
	}
	if v := os.Getenv("PORT"); v != "" {
		return ":" + v
	}

	return ":8080"
}

func firstEnv(keys ...string) string {
	for _, key := range keys {
		if v := os.Getenv(key); v != "" {
			return v
		}
	}

	return ""
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
