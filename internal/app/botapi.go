package app

import (
	"regexp"
	"strings"

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
		message := strings.ReplaceAll(update.Message.Text, " ", "")
		if update.Message.Entities != nil && update.Message.ForwardFromChat != nil &&
			b.isValidationUser(update.Message.ForwardFromChat.UserName) &&
			b.isValidationMessage(message) {
			continue
		} else if update.Message.Entities != nil && update.Message.From != nil &&
			b.isValidationUser(update.Message.From.UserName) &&
			b.isValidationMessage(message) {
			continue
		}
		delMessageConfig := tgbotapi.DeleteMessageConfig{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.MessageID,
		}
		b.bot.DeleteMessage(delMessageConfig)
	}
	return nil
}

func (b *BotAPI) isValidationUser(username string) bool {
	for _, accessUsername := range b.config.Access {
		if accessUsername == username {
			return true
		}
	}
	return false
}

func (b *BotAPI) isValidationMessage(message string) bool {
	re, err := regexp.Compile(b.config.Regexp)
	if err != nil {
		return false
	}
	if !re.Match([]byte(message)) {
		return false
	}
	return true
}
