package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Addr     string `yaml:"addr"`
	DBDriver string `yaml:"db_driver"`
	DBDSN    string `yaml:"db_dsn"`
}

var Conf Config

func MustInit() {
		wd, _ := os.Getwd()

	def := filepath.Join(wd, "internal", "config", "config.yaml")
	path := getenv("CONFIG_PATH", def)
	c, err := fromYAML(path)
	if err != nil {
		panic("cannot load config from path. error : "+err.Error())
	}
 Conf = c
}

func fromYAML(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var c Config
	if err := yaml.Unmarshal(b, &c); err != nil {
		return Config{}, err
	}
	applyDefaults(&c)
	return c, nil
}

func FromEnv() Config {
	addr := getenv("ADDR", ":8080")
	driver := getenv("DB_DRIVER", "sqlite")
	dsn := getenv("DB_DSN", "file:memdb1?mode=memory&cache=shared")
	return Config{Addr: addr, DBDriver: driver, DBDSN: dsn}
}

func getenv(k, def string) string {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	return v
}

func applyDefaults(c *Config) {
	if c.Addr == "" {
		c.Addr = ":8080"
	}
	if c.DBDriver == "" {
		c.DBDriver = "sqlite"
	}
	if c.DBDSN == "" {
		c.DBDSN = "file:memdb1?mode=memory&cache=shared"
	}
}
