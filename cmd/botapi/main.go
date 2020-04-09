package main

import (
	"flag"
	"log"

	"github.com/b2r2/chat-controller-bot/internal/app"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "./configs/config.toml", "path to config file")
}

func main() {
	flag.Parse()
	config := app.NewConfig()
	if _, err := toml.DecodeFile(configPath, config); err != nil {
		log.Fatal(err)
	}
	b := app.NewBotAPI(config)
	if err := b.Start(); err != nil {
		log.Fatal(err)
	}
}
