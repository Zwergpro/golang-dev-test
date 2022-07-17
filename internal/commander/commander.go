package commander

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"log"
)

type Commander struct {
	bot *tgbotapi.BotAPI
}

func Init(tgApiKey string) (*Commander, error) {
	bot, err := tgbotapi.NewBotAPI(tgApiKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create bot")
	}

	bot.Debug = true

	log.Printf("Authorized on account: %s", bot.Self.UserName)

	return &Commander{
		bot: bot,
	}, nil
}

func (c *Commander) Run() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		_, err := c.bot.Send(msg)
		if err != nil {
			return errors.Wrap(err, "failed to send message")
		}

	}

	return nil
}
