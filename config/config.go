package config

import (
	"net"
	"net/url"
	"os"
)

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
}

func (d DBConfig) DBUrl() string {
	u := url.URL{
		Scheme: d.User,
		User: url.UserPassword(d.User, d.Password),
		Host: net.JoinHostPort(d.Host, d.Port),
		Path: "/" + d.Name,
	}

	query := u.Query()
	query.Set("sslmode", d.SSLMode)
	u.RawQuery = query.Encode()

	return u.String()
}
