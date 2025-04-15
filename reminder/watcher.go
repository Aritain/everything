package reminder

import (
	"fmt"
	"time"

	"everything/common"
	"everything/models"
	r "everything/models/reminder"
)

const TIMEOUT = 30

func WatchReminders(config *models.Config) {
	location, _ := time.LoadLocation(config.TimezoneLocation)
	for {
		reminders := LoadReminders()
		now := time.Now()
		for _, reminder := range reminders {
			reminderTime := reminder.ReminderData.NextReminder.In(location)
			if now.After(reminderTime) {
				SendReminder(reminder.ReminderData)
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

func SendReminder(reminder r.Reminder) {
	msgText := fmt.Sprintf("It's time for *%s*.", reminder.ReminderText)
	userID := reminder.UserID
	go common.SendTGMessage(userID, msgText)
}
