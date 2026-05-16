package config

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type Config struct {
	Server ServerConfig `mapstructure:"server"`
}

func (c *Config) GetServerPort() int {
	return c.Server.Port
}

func Default() *Config {
	return &Config{
		Server: ServerConfig{
			Port: 8080,
		},
	}
}
