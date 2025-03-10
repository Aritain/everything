package reminder

import (
	"fmt"
	"slices"
	"strconv"

	"everything/models"
	r "everything/models/reminder"
)

func ReadReminderName(ri *r.ReminderInput) (mr models.ModuleResponse) {
	if len(ri.Text) > 30 {
		mr.ResponseText = "Too long"
		mr.ResponseCode = true
		return mr
	}
	AppendCache(ri.ReminderCache, ri.UserID, r.Reminder{ReminderText: ri.Text})
	mr.ResponseText = "When?"
	return mr
}

func ReadReminderTime(ri *r.ReminderInput) (mr models.ModuleResponse) {
	reminderTime, err := ParseTime(ri.Text)
	if err != nil {
		mr.ResponseText = "Failed to read time provided"
		mr.ResponseCode = true
		return mr
	}
	AppendCache(ri.ReminderCache, ri.UserID, r.Reminder{NextReminder: reminderTime})
	mr.ResponseText = "Repeat?"
	return mr
}

func ReadReminderRepeat(ri *r.ReminderInput) (mr models.ModuleResponse) {
	if ri.Text == "No" {
		mr.ResponseText = "Done"
		mr.ResponseCode = true
		WriteReminder(ri.ReminderCache, ri.UserID)
		DeleteReminderCache(ri.ReminderCache, ri.UserID)
		return mr
	}
	AppendCache(ri.ReminderCache, ri.UserID, r.Reminder{RepeatToggle: true})
	mr.ResponseText = "Mode?"
	return mr
}

func ReadReminderMode(ri *r.ReminderInput) (mr models.ModuleResponse) {
	allowedResponces := []string{"day", "week", "month", "year"}
	if !slices.Contains(allowedResponces, ri.Text) {
		mr.ResponseText = "day/week/month/year"
		mr.ResponseCode = true
		return mr
	}
	AppendCache(ri.ReminderCache, ri.UserID, r.Reminder{RepeatMode: ri.Text})
	mr.ResponseText = "Value?"
	return mr
}

func ReadReminderValue(ri *r.ReminderInput) (mr models.ModuleResponse) {
	value64, err := strconv.ParseUint(ri.Text, 10, 8)
	if err != nil {
		mr.ResponseText = "Bad value"
		mr.ResponseCode = true
		return mr
	}
	value8 := uint8(value64)
	AppendCache(ri.ReminderCache, ri.UserID, r.Reminder{RepeatValue: value8})
	mr.ResponseText = "Done"
	WriteReminder(ri.ReminderCache, ri.UserID)
	DeleteReminderCache(ri.ReminderCache, ri.UserID)
	return mr
}

func GetReminders(userID int64) (mr models.ModuleResponse) {
	reminders := LoadReminders()
	for _, reminder := range reminders {
		if reminder.ReminderData.UserID == userID {
			mr.ResponseText += FormatReminder(reminder.ReminderData)
		}
	}
	if len(mr.ResponseText) == 0 {
		mr.ResponseText = "No reminders found."
	}
	return mr
}

func DeleteReminderQuery(userID int64) (mr models.ModuleResponse) {
	var counter int
	response := "Send me the number of reminder to delete:\n"
	mr.ResponseText += response
	reminders := LoadReminders()
	for _, reminder := range reminders {
		if reminder.ReminderData.UserID == userID {
			counter += 1
			mr.ResponseText += fmt.Sprintf("(%v) ", counter)
			mr.ResponseText += FormatReminder(reminder.ReminderData)
		}
	}
	if mr.ResponseText == response {
		mr.ResponseText = "No reminders found."
	}
	return mr
}

func DeleteReminderConfirm(input string, userID int64) (mr models.ModuleResponse) {
	number, err := strconv.Atoi(input)
	if (err != nil) || (number <= 0) {
		mr.ResponseText = "Bad value"
		mr.ResponseCode = true
		return mr
	}
	// Reduce number for it to match index
	number -= 1
	reminders := LoadReminders()
	var userReminders []r.ReminderFile
	for _, reminder := range reminders {
		if reminder.ReminderData.UserID == userID {
			userReminders = append(userReminders, reminder)
		}
	}
	if number > len(userReminders) {
		mr.ResponseText = "Bad value"
		mr.ResponseCode = true
		return mr
	}
	if !DeleteReminder(userReminders[number].FileName) {
		mr.ResponseText = "Failed to delete the file"
		mr.ResponseCode = true
		return mr
	}
	mr.ResponseText = "Done"
	return mr
}
