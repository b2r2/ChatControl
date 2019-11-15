package app

// Config ...
type Config struct {
	Token     string   `toml:"token"`
	DebugMode bool     `toml:"debug_mode"`
	LogLevel  string   `toml:"log_level"`
	Regexp    string   `toml:"regexp"`
	Access    []string `toml:"access_usernames"`
}

// NewConfig return struct config
func NewConfig() *Config {
	return &Config{
		DebugMode: true,
		LogLevel:  "debug",
	}
}
