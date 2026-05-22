package cfg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"friday/pkg/cfg"
)

func TestLoad_defaults(t *testing.T) {
	for _, key := range []string{"HTTP_ADDR", "PORT", "DATABASE_URL", "POSTGRES_URL", "POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB", "RESEND_API_KEY", "RESEND_FROM", "SMTP_HOST", "SMTP_PORT", "SMTP_USERNAME", "SMTP_PASSWORD", "SMTP_FROM"} {
		t.Setenv(key, "")
	}

	c := cfg.Load()

	assert.Equal(t, ":8080", c.HTTPAddr)
	assert.Equal(t, "localhost", c.Postgres.Host)
	assert.Equal(t, "5432", c.Postgres.Port)
	assert.Equal(t, "postgres", c.Postgres.User)
	assert.Equal(t, "postgres", c.Postgres.Password)
	assert.Equal(t, "friday", c.Postgres.Database)
	assert.Equal(t, "", c.Resend.APIKey)
	assert.Equal(t, "onboarding@resend.dev", c.Resend.From)
	assert.Equal(t, "", c.SMTP.Host)
	assert.Equal(t, "587", c.SMTP.Port)
	assert.Equal(t, "", c.SMTP.Username)
	assert.Equal(t, "", c.SMTP.Password)
	assert.Equal(t, "", c.SMTP.From)
}

func TestLoad_envOverride(t *testing.T) {
	t.Setenv("HTTP_ADDR", ":9090")
	t.Setenv("POSTGRES_HOST", "db.prod.internal")
	t.Setenv("POSTGRES_PORT", "5433")
	t.Setenv("POSTGRES_USER", "admin")
	t.Setenv("POSTGRES_PASSWORD", "s3cr3t")
	t.Setenv("POSTGRES_DB", "myapp")
	t.Setenv("RESEND_API_KEY", "re_test")
	t.Setenv("RESEND_FROM", "Friday <auth@example.com>")
	t.Setenv("SMTP_HOST", "smtp.example.com")
	t.Setenv("SMTP_PORT", "2525")
	t.Setenv("SMTP_USERNAME", "auth@example.com")
	t.Setenv("SMTP_PASSWORD", "smtp-secret")
	t.Setenv("SMTP_FROM", "Friday <auth@example.com>")

	c := cfg.Load()

	assert.Equal(t, ":9090", c.HTTPAddr)
	assert.Equal(t, "db.prod.internal", c.Postgres.Host)
	assert.Equal(t, "5433", c.Postgres.Port)
	assert.Equal(t, "admin", c.Postgres.User)
	assert.Equal(t, "s3cr3t", c.Postgres.Password)
	assert.Equal(t, "myapp", c.Postgres.Database)
	assert.Equal(t, "re_test", c.Resend.APIKey)
	assert.Equal(t, "Friday <auth@example.com>", c.Resend.From)
	assert.Equal(t, "smtp.example.com", c.SMTP.Host)
	assert.Equal(t, "2525", c.SMTP.Port)
	assert.Equal(t, "auth@example.com", c.SMTP.Username)
	assert.Equal(t, "smtp-secret", c.SMTP.Password)
	assert.Equal(t, "Friday <auth@example.com>", c.SMTP.From)
}
