package postgres

import (
	"fmt"
	"os"
)

type Config struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

func LoadDBConfig() (Config, error) {
	getEnv := func(key string) (string, error) {
		value, ok := os.LookupEnv(key)
		if !ok {
			return "", fmt.Errorf("environment variable %s not found", key)
		}
		return value, nil
	}

	host, err := getEnv("DB_HOST")
	if err != nil {
		return Config{}, err
	}

	port, err := getEnv("DB_PORT")
	if err != nil {
		return Config{}, err
	}

	name, err := getEnv("DB_NAME")
	if err != nil {
		return Config{}, err
	}

	user, err := getEnv("DB_USER")
	if err != nil {
		return Config{}, err
	}

	password, err := getEnv("DB_PASS")
	if err != nil {
		return Config{}, err
	}

	return Config{
		Host:     host,
		Port:     port,
		Name:     name,
		User:     user,
		Password: password,
	}, nil
}

func (c Config) ConnString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host,
		c.Port,
		c.User,
		c.Password,
		c.Name,
	)
}
