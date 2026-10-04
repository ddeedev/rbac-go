package config

import (
	"fmt"
	"net/url"
)

// build connection uri w/o set monog user/pass auth later
func BuildMongoURI(cfg *DatabaseConfig) (string, error) {
	host := cfg.Host
	if host == "" {
		return "", fmt.Errorf("MONGO_DB_HOST is not set")
	}

	port := cfg.Port
	user := cfg.User
	password := cfg.Password
	// incase no need auth
	if user == "" || password == "" {
		return fmt.Sprintf("mongodb://%s:%d", host, port), nil
	}

	uri := &url.URL{
		Scheme: "mongodb",
		User:   url.UserPassword(user, password),
		Host:   fmt.Sprintf("%s:%d", host, port),
	}

	return uri.String(), nil
}
