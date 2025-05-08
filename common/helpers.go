package common

import (
	"log"
	"time"

	"everything/models"

	t "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const TIMEOUT = 10

func SendTGMessage(tgm models.TGMessage) {
	bot, _ := t.NewBotAPI(tgm.TGToken)
	msg := t.NewMessage(tgm.UserID, tgm.Text)
	if tgm.ParseMode == "" {
		tgm.ParseMode = "Markdown"
	}
	if len(tgm.Keyboard.InlineKeyboard) == 0 {
		tgm.Keyboard = t.InlineKeyboardMarkup{InlineKeyboard: make([][]t.InlineKeyboardButton, 0)}
	}
	msg.ParseMode = tgm.ParseMode
	msg.ReplyMarkup = tgm.Keyboard
	var err error

	for {
		_, err = bot.Send(msg)
		if err == nil {
			break
		}
		log.Print(err)
		time.Sleep(TIMEOUT * time.Second)
	}
}
