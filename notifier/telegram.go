package notifier

import (
	"fmt"

	"github.com/idoberko2/home_health_be/general"
	"github.com/pkg/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewTelegram(cfg TelegramConfig) Notifier {
	return &telegramNotifier{
		cfg: cfg,
	}
}

type telegramNotifier struct {
	cfg TelegramConfig
	bot *tgbotapi.BotAPI
}

func (t *telegramNotifier) Init() error {
	bot, err := tgbotapi.NewBotAPI(t.cfg.Token)
	if err != nil {
		return err
	}

	if t.cfg.IsDebug {
		bot.Debug = true
	}

	t.bot = bot

	return nil
}

func (t *telegramNotifier) NotifyStateChange(state general.State) error {
	if t.bot == nil {
		return ErrNotInitialized
	}

	msg := tgbotapi.NewMessage(t.cfg.ChatID, getMessage(state))
	_, err := t.bot.Send(msg)

	return errors.Wrap(err, "error sending message")
}

func getMessage(state general.State) string {
	return fmt.Sprintf("State changed to '%s'", state)
}

var ErrNotInitialized = errors.New("not initialized")
