package app

// Config ...
type Config struct {
	Token          string   `toml:"token"`
	DebugMode      bool     `toml:"debug_mode"`
	LogLevel       string   `toml:"log_level"`
	Regexp         string   `toml:"regexp"`
	AccessUsers    []string `toml:"access_users"`
	AccessChannels []string `toml:"access_channels"`
}

// NewConfig return struct config
func NewConfig() *Config {
	return &Config{
		DebugMode: true,
		LogLevel:  "debug",
	}
}
