package configuration

import "log/slog"

type Configuration struct {
	DBConf         DatabaseConfiguration
	Env            string
	GameServerPort uint16
	LogLevel       slog.Level
}

type DatabaseConfiguration struct {
	Host     string
	Password string
	Username string
}

func GetConfiguration() *Configuration {
	switch env := GetEnv(); env {
	case DevEnv:
		return &Configuration{
			Env:            env,
			GameServerPort: 27015,
			LogLevel:       slog.LevelDebug,
		}
	}

	return &Configuration{}
}
