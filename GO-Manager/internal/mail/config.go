package mail

import (
	"fmt"
	"os"
)

type Config struct {
	Host string
	Port string
	User string
	Pass string
	From string
}

func LoadConfig() Config {
	return Config{
		Host: os.Getenv("SMTP_HOST"),
		Port: os.Getenv("SMTP_PORT"),
		User: os.Getenv("SMTP_USER"),
		Pass: os.Getenv("SMTP_PASS"),
		From: os.Getenv("SMTP_FROM"),
	}
}

func (c Config) Ready() error {
	if c.Host == "" || c.Port == "" || c.From == "" {
		return fmt.Errorf("smtp not configured")
	}
	return nil
}
