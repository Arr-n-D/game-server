package configuration

type Configuration struct {
	DBConf         DatabaseConfiguration
	SentryURL      string
	Env            string
	GameServerPort uint16
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
			// sentry.Dsn
		}
	}

	return &Configuration{}

}
