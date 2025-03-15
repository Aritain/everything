package reminder

import (
	"strings"
	"time"

	c "everything/config"
	r "everything/models/reminder"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func WatchReminders() {
	for {
		reminders := LoadReminders()
		now := time.Now()
		for _, reminder := range reminders {
			if now.After(reminder.ReminderData.NextReminder) {
				go SendReminder(reminder.ReminderData)
				DeleteReminder(reminder.FileName)
				if reminder.ReminderData.RepeatToggle {
					bumpReminder := reminder.ReminderData
					bumpReminder.NextReminder = UpdateReminder(bumpReminder)
					PrepareReminderWrite(bumpReminder)
				}
			}
		}
		time.Sleep(30 * time.Second)
	}
}

func PrepareReminderWrite(reminder r.Reminder) {
	var reminders []r.Reminder
	reminders = append(reminders, reminder)
	WriteReminder(&reminders, reminder.UserID)
}

func SendReminder(reminder r.Reminder) {
	msgText := "It's time for "
	msgText += FormatReminder(reminder)
	msgText = strings.Replace(msgText, "for Reminder", "for", -1)
	userID := reminder.UserID
	config, _ := c.LoadConfig()
	bot, _ := tgbotapi.NewBotAPI(config.TGToken)
	msg := tgbotapi.NewMessage(userID, msgText)
	var err error
	for {
		_, err = bot.Send(msg)
		if err == nil {
			break
		}
	}
}
