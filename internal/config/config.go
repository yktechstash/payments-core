package config

import (
	"os"
)

type Config struct {
	Addr     string // e.g. ":8080"
	DBDriver string // "postgres" or "sqlite"
	DBDSN    string // connection string / DSN
}

func FromEnv() Config {
	addr := getenv("ADDR", ":8080")
	driver := getenv("DB_DRIVER", "postgres")
	dsn := getenv("DB_DSN", "postgres://postgres:postgres@localhost:5432/aaa?sslmode=disable")
	return Config{
		Addr:     addr,
		DBDriver: driver,
		DBDSN:    dsn,
	}
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}
