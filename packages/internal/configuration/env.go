package configuration

import (
	"os"
)

const (
	DevEnv   = "dev"
	ProdEnv  = "prod"
	LocalEnv = "local"
)

func GetEnv() string {
	env := os.Getenv("ENV")

	if env == "" {
		return ""
	}

	return env
}

func (conf *Configuration) IsDevEnv() bool {
	return conf.Env == DevEnv
}

func (conf *Configuration) IsLocalEnv() bool {
	return conf.Env == LocalEnv
}
