package configuration

import (
	"os"
)

const (
	DevEnv  = "dev"
	ProdEnv = "prod"
)

func GetEnv() string {
	env := os.Getenv("ENV")

	if env == "" {
		panic("No environment found")
	}

	return env
}
