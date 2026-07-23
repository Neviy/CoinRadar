// package bot содержит логику для работы с Telegram-ботом.
package bot

import (
	"CoinRadar/internal/services"
	"time"

	tele "gopkg.in/telebot.v4"
)

// Bot представляет Telegram-бота и его зависимости.
type Bot struct {
	bot *tele.Bot

	userService         *services.UserService
	coinService         *services.CoinService
	subscriptionService *services.SubscriptionService
}

// NewBot создает новый экземпляр Telegram-бота.
func NewBot(token string, userService *services.UserService, coinService *services.CoinService,
	subscriptionService *services.SubscriptionService) (*Bot, error) {

	tgBot, err := tele.NewBot(tele.Settings{
		Token: token,
		Poller: &tele.LongPoller{
			Timeout: 10 * time.Second,
		},
	})
	if err != nil {
		return nil, err
	}

	return &Bot{
		bot:                 tgBot,
		userService:         userService,
		coinService:         coinService,
		subscriptionService: subscriptionService,
	}, nil
}

// Start запускает Telegram-бота.
func (b *Bot) Start() {
	b.bot.Start()
}

// Stop останавливает Telegram-бота.
func (b *Bot) Stop() {
	b.bot.Stop()
}
