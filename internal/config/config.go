package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Domain      string
	APILogin    string
	APIKey      string
	PVEName     string
	TLSInsecure bool
}

func LoadENV() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	tlsInsecure, _ := strconv.ParseBool(os.Getenv("TLS_INSECURE"))

	return &Config{
		Domain:      os.Getenv("DOMAIN"),
		APILogin:    os.Getenv("APILOGIN"),
		APIKey:      os.Getenv("APIKEY"),
		PVEName:     os.Getenv("PVENAME"),
		TLSInsecure: tlsInsecure,
	}, nil
}
