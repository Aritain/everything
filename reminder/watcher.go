package reminder

import (
	"strings"
	"time"

	"everything/common"
	r "everything/models/reminder"
)

const TIMEOUT = 30

func WatchReminders() {
	for {
		reminders := LoadReminders()
		now := time.Now()
		for _, reminder := range reminders {
			if now.After(reminder.ReminderData.NextReminder) {
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
	msgText := "It's time for "
	msgText += FormatReminder(reminder)
	msgText = strings.Replace(msgText, "for Reminder", "for", -1)
	userID := reminder.UserID
	go common.SendTGMessage(userID, msgText)
}
