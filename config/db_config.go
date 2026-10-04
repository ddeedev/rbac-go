package config

import (
	"fmt"
	"net/url"
	"os"
)

// build connection uri w/o set monog user/pass auth later
func BuildMongoURI() (string, error) {
	host := os.Getenv("MONGO_DB_HOST")
	if host == "" {
		return "", fmt.Errorf("MONGO_DB_HOST is not set")
	}

	port := os.Getenv("MONGO_DB_PORT")
	if port == "" {
		return "", fmt.Errorf("MONGO_DB_PORT is not set")
	}

	user := os.Getenv("MONGO_DB_USERNAME")
	password := os.Getenv("MONGO_DB_PASSWORD")

	// incase no need auth
	if user == "" || password == "" {
		return fmt.Sprintf("mongodb://%s:%s", host, port), nil
	}

	uri := &url.URL{
		Scheme: "mongodb",
		User:   url.UserPassword(user, password),
		Host:   fmt.Sprintf("%s:%s", host, port),
	}

	return uri.String(), nil
}
