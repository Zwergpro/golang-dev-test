package commander

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"homework-1/internal/repository"
	"log"
)

type CmdHandler func(repository.Product, string) string

type Commander struct {
	bot               *tgbotapi.BotAPI
	router            map[string]CmdHandler
	ProductRepository repository.Product
}

func Init(tgApiKey string, repository repository.Product) (*Commander, error) {
	bot, err := tgbotapi.NewBotAPI(tgApiKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create bot")
	}

	bot.Debug = true

	log.Printf("Authorized on account: %s", bot.Self.UserName)

	return &Commander{
		bot:               bot,
		router:            make(map[string]CmdHandler),
		ProductRepository: repository,
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
		if update.Message.Command() != "" {
			if cmd, ok := c.router[update.Message.Command()]; ok {
				msg.Text = cmd(c.ProductRepository, update.Message.CommandArguments())
			} else {
				msg.Text = fmt.Sprintf("Invalid command: %v", update.Message.Command())
			}
		}

		_, err := c.bot.Send(msg)
		if err != nil {
			return errors.Wrap(err, "failed to send message")
		}

	}

	return nil
}

func (c *Commander) RegisterHandler(cmd string, handler CmdHandler) {
	c.router[cmd] = handler
}
