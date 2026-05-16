package config

type LoggingConfig struct {
	Level  string
	Format string
}

type ServerConfig struct {
	Port              int
	AccessLogsEnabled bool
}

type Config struct {
	Server  ServerConfig
	Logging LoggingConfig
}

func (c *Config) GetServerPort() int {
	return c.Server.Port
}

func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Port:              8080,
			AccessLogsEnabled: true,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
	}
}
