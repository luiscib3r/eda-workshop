package telegram

import (
	"backend/internal/infrastructure/config"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot    *tgbotapi.BotAPI
	chatId int64
}

func NewTelegramBot(
	cfg *config.AppConfig,
) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.BotToken)
	if err != nil {
		return nil, err
	}

	return &TelegramBot{
		bot:    bot,
		chatId: cfg.Telegram.ChatID,
	}, nil
}

func (t *TelegramBot) SendMessage(message string) error {
	msg := tgbotapi.NewMessage(t.chatId, escapeMarkdownV2(message))
	msg.ParseMode = "MarkdownV2"
	_, err := t.bot.Send(msg)
	return err
}

func (t *TelegramBot) Start() {
	cfg := tgbotapi.NewUpdate(0)
	cfg.Timeout = 30

	updates := t.bot.GetUpdatesChan(cfg)

	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			// If is command /ping, respond with "pong {chatID}"
			if update.Message.IsCommand() && update.Message.Command() == "ping" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("pong %d", update.Message.Chat.ID))
				t.bot.Send(msg)
			}
		}
	}()
}

var markdownV2Replacer = strings.NewReplacer(
	"\\", "\\\\",
	"_", "\\_",
	"[", "\\[",
	"]", "\\]",
	"(", "\\(",
	")", "\\)",
	"~", "\\~",
	"`", "\\`",
	">", "\\>",
	"#", "\\#",
	"+", "\\+",
	"-", "\\-",
	"=", "\\=",
	"|", "\\|",
	"{", "\\{",
	"}", "\\}",
	".", "\\.",
	"!", "\\!",
)

func escapeMarkdownV2(text string) string {
	return markdownV2Replacer.Replace(text)
}
