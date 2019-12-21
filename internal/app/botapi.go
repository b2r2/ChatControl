package app

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

// BotAPI ...
type BotAPI struct {
	config  *Config
	logger  *logrus.Logger
	bot     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
}

// NewBotAPI ...
func NewBotAPI(config *Config) *BotAPI {
	return &BotAPI{
		config: config,
		logger: logrus.New(),
	}
}

// Start ...
func (b *BotAPI) Start() error {
	if err := b.configureLogger(); err != nil {
		return err
	}

	if err := b.configureBot(); err != nil {
		return err
	}

	b.logger.Info("Starting bot. Debug mode:", b.bot.Debug)
	if err := b.Handler(); err != nil {
		return err
	}
	return nil
}

func (b *BotAPI) configureBot() error {
	bot, err := tgbotapi.NewBotAPI(b.config.Token)
	if err != nil {
		return err
	}

	b.bot = bot
	b.bot.Debug = b.config.DebugMode

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}
	b.updates = updates
	return nil
}
func (b *BotAPI) configureLogger() error {
	level, err := logrus.ParseLevel(b.config.LogLevel)
	if err != nil {
		return err
	}
	b.logger.SetLevel(level)
	return nil
}

// Handler ...
func (b *BotAPI) Handler() error {
	for update := range b.updates {
		if update.Message == nil {
			continue
		}
		deleteMessageConfig := tgbotapi.DeleteMessageConfig{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.MessageID,
		}
		if update.Message.Entities != nil && update.Message.ForwardFromChat != nil && !b.isValidationChannel(update.Message.ForwardFromChat.UserName) {
			b.bot.DeleteMessage(deleteMessageConfig)
		} else if update.Message.Entities != nil && update.Message.From != nil && !b.isValidationUser(update.Message.From.UserName) {
			b.bot.DeleteMessage(deleteMessageConfig)
		}
	}
	return nil
}

func (b *BotAPI) isValidationChannel(c string) bool {
	for _, channel := range b.config.AccessChannels {
		if channel == c {
			return true
		}
	}
	return false
}

func (b *BotAPI) isValidationUser(u string) bool {
	for _, accessusers := range b.config.AccessUsers {
		if u == accessusers {
			return true
		}
	}
	return false
}
