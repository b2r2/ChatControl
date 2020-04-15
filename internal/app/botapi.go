package app

import (
	"regexp"
	"strings"
	"unicode"

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
		b.logger.Printf("configureBot: %s\n", err)
		return err
	}

	b.logger.Info("Starting bot. Debug mode:", b.bot.Debug)
	if err := b.handler(); err != nil {
		b.logger.Printf("handler: %s\n", err)
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

func (b *BotAPI) handler() error {
	for update := range b.updates {
		if update.Message == nil {
			continue
		}
		mc := tgbotapi.DeleteMessageConfig{
			ChatID:    update.Message.Chat.ID,
			MessageID: update.Message.MessageID,
		}
		isSticker := b.config.StickerMode && update.Message.Sticker != nil
		isChannel := update.Message.ForwardFromChat != nil && !b.isVerifyChannel(update.Message.ForwardFromChat.UserName)
		isUser := update.Message.From != nil && !b.isVerifyUser(update.Message.From.UserName)
		isLink := update.Message.From != nil && b.isContainedLink(update.Message.Text)
		isEntities := update.Message.Entities != nil || update.Message.CaptionEntities != nil

		isRemoveMessage := (isSticker || isEntities || isLink) && (isChannel || isUser)
		if isRemoveMessage {
			if _, err := b.bot.DeleteMessage(mc); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *BotAPI) isVerifyChannel(c string) bool {
	for _, channel := range b.config.AccessChannels {
		if channel == c {
			return true
		}
	}
	return false
}

func (b *BotAPI) isVerifyUser(u string) bool {
	for _, user := range b.config.AccessUsers {
		if u == user {
			return true
		}
	}
	return false
}

func (b *BotAPI) isContainedLink(s string) bool {
	return b.handleRegexp(b.handleText(s)) //or Link Validation?
}

func (b *BotAPI) handleRegexp(s string) bool {
	return regexp.MustCompile(b.config.Regexp).MatchString(s)
}

func (b *BotAPI) handleText(s string) string {
	return strings.Join(strings.FieldsFunc(s, func(c rune) bool {
		return unicode.IsSpace(c)
	}), "")
}
