package common

import (
	c "everything/config"

	t "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendTGMessage(userID int64, text string) {
	config, _ := c.LoadConfig()
	bot, _ := t.NewBotAPI(config.TGToken)
	msg := t.NewMessage(userID, text)
	msg.ParseMode = "Markdown"
	var err error

	for {
		_, err = bot.Send(msg)
		if err == nil {
			break
		}
	}
}
