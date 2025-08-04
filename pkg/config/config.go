package config

import (
	"flag"
	"os"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// getEnv reads an environment variable or returns the default value
func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}

var (
	dbhost   = flag.String("host", getEnv("DB_HOST", "localhost"), "host name")
	port     = flag.String("port", getEnv("DB_PORT", "5432"), "port number")
	user     = flag.String("user", getEnv("DB_USER", "postgres"), "user name")
	password = flag.String("password", getEnv("DB_PASSWORD", "postgres"), "password")
	dbname   = flag.String("dbname", getEnv("DB_NAME", "onset_adaptar"), "database name")
)

func LoadDBConfig() *DBConfig {
	flag.Parse()

	return &DBConfig{
		Host:     *dbhost,
		Port:     *port,
		User:     *user,
		Password: *password,
		DBName:   *dbname,
	}
}
