package cfg_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"friday/pkg/cfg"
)

func TestLoad_defaults(t *testing.T) {
	for _, key := range []string{"HTTP_ADDR", "POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB"} {
		t.Setenv(key, "")
	}

	c := cfg.Load()

	assert.Equal(t, ":8080", c.HTTPAddr)
	assert.Equal(t, "localhost", c.Postgres.Host)
	assert.Equal(t, "5432", c.Postgres.Port)
	assert.Equal(t, "postgres", c.Postgres.User)
	assert.Equal(t, "postgres", c.Postgres.Password)
	assert.Equal(t, "friday", c.Postgres.Database)
}

func TestLoad_envOverride(t *testing.T) {
	t.Setenv("HTTP_ADDR", ":9090")
	t.Setenv("POSTGRES_HOST", "db.prod.internal")
	t.Setenv("POSTGRES_PORT", "5433")
	t.Setenv("POSTGRES_USER", "admin")
	t.Setenv("POSTGRES_PASSWORD", "s3cr3t")
	t.Setenv("POSTGRES_DB", "myapp")

	c := cfg.Load()

	assert.Equal(t, ":9090", c.HTTPAddr)
	assert.Equal(t, "db.prod.internal", c.Postgres.Host)
	assert.Equal(t, "5433", c.Postgres.Port)
	assert.Equal(t, "admin", c.Postgres.User)
	assert.Equal(t, "s3cr3t", c.Postgres.Password)
	assert.Equal(t, "myapp", c.Postgres.Database)
}
