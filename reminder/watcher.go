package reminder

import (
	"fmt"
	"time"

	"everything/common"
	cfg "everything/config"
	"everything/models"
	r "everything/models/reminder"
)

const TIMEOUT = 30

func WatchReminders() {
	config := cfg.Get().Config()
	location, _ := time.LoadLocation(config.TimezoneLocation)
	for {
		reminders := LoadReminders()
		now := time.Now()
		for _, reminder := range reminders {
			reminderTime := reminder.ReminderData.NextReminder.In(location)
			if now.After(reminderTime) {
				SendReminder(config.TGToken, reminder.ReminderData)
				DeleteReminder(reminder.FileName)
				if reminder.ReminderData.RepeatToggle {
					bumpReminder := reminder.ReminderData
					bumpReminder.NextReminder = UpdateReminder(bumpReminder)
					PrepareReminderWrite(bumpReminder)
				}
			}
		}
		time.Sleep(TIMEOUT * time.Second)
	}
}

func PrepareReminderWrite(reminder r.Reminder) {
	var reminders []r.Reminder
	reminders = append(reminders, reminder)
	WriteReminder(&reminders, reminder.UserID)
}

func SendReminder(tgToken string, reminder r.Reminder) {
	msgText := fmt.Sprintf("It's time for *%s*.", reminder.ReminderText)
	var tgm models.TGMessage
	tgm.TGToken = tgToken
	tgm.UserID = reminder.UserID
	tgm.Text = msgText
	go common.SendTGMessage(tgm)
}
