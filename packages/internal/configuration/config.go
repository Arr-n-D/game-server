package configuration

import (
	"log/slog"
	"os"
)

type databaseSecrets struct {
}

type Configuration struct {
	Env             string
	GameServerPort  uint16
	LogLevel        slog.Level
	AwsRegion       string
	DatabaseSecrets databaseSecrets
}

func GetConfiguration() *Configuration {
	switch env := GetEnv(); env {
	case DevEnv:
		return &Configuration{
			Env:            env,
			GameServerPort: 27015,
			LogLevel:       slog.LevelDebug,
			AwsRegion:      os.Getenv("AWS_REGION"),
		}
	case LocalEnv:
		return &Configuration{
			Env:            env,
			GameServerPort: 27015,
			LogLevel:       slog.LevelDebug,
			AwsRegion:      os.Getenv("AWS_REGION"),
		}
	}

	return &Configuration{}
}

func (conf *Configuration) GetAwsRegion() string {
	return conf.AwsRegion
}
