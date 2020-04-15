package app

// Config ...
type Config struct {
	Token          string   `toml:"token"`
	DebugMode      bool     `toml:"debug_mode"`
	LogLevel       string   `toml:"log_level"`
	AccessUsers    []string `toml:"access_users"`
	AccessChannels []string `toml:"access_channels"`
	StickerMode    bool     `toml:"sticker_mode"`
	Regexp         string   `toml:"regexp"`
}

// NewConfig return struct config
func NewConfig() *Config {
	return &Config{
		DebugMode: true,
		LogLevel:  "debug",
	}
}
