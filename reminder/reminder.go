package reminder

import (
	"fmt"
	"log"
	"slices"
	"strconv"

	"everything/common"
	"everything/models"
	r "everything/models/reminder"
)

func ReminderCreationStart(userID int64, rc *[]r.Reminder) (mr models.ModuleResponse) {
	*rc = append(*rc, r.Reminder{UserID: userID})
	mr.Text = "Reminder name?"
	mr.Keyboard = common.CompileCancelKeyboard()
	return mr
}

func ReadReminderName(ri *r.ReminderInput) (mr models.ModuleResponse) {
	mr.Keyboard = common.CompileCancelKeyboard()
	if len(ri.Text) > 30 {
		mr.Text = "Too long"
		mr.Error = true
		return mr
	}
	AppendCache(ri.ReminderCache, ri.UserID, r.Reminder{ReminderText: ri.Text})
	mr.Text = "When?"
	return mr
}

func ReadReminderTime(ri *r.ReminderInput, config *models.Config) (mr models.ModuleResponse) {
	reminderTime, err := ParseTime(ri.Text, config)
	if err != nil {
		mr.Text = "Failed to read time provided"
		mr.Error = true
		return mr
	}
	AppendCache(ri.ReminderCache, ri.UserID, r.Reminder{NextReminder: reminderTime})
	mr.Text = "Repeat?"
	mr.Keyboard = common.CompileYesNoKeyboard()
	return mr
}

func ReadReminderRepeat(ri *r.ReminderInput) (mr models.ModuleResponse) {
	if ri.Text == "No" {
		mr.Text = "Done"
		mr.EndChat = true
		mr.Keyboard = common.CompileDefaultKeyboard()
		WriteReminder(ri.ReminderCache, ri.UserID)
		DeleteReminderCache(ri.ReminderCache, ri.UserID)
		return mr
	}
	AppendCache(ri.ReminderCache, ri.UserID, r.Reminder{RepeatToggle: true})
	mr.Text = "Mode?"
	mr.Keyboard = common.CompileReminderModeKeyboard()
	return mr
}

func ReadReminderMode(ri *r.ReminderInput) (mr models.ModuleResponse) {
	allowedResponses := []string{"day", "week", "month", "year"}
	if !slices.Contains(allowedResponses, ri.Text) {
		mr.Text = "day/week/month/year"
		mr.Error = true
		return mr
	}
	AppendCache(ri.ReminderCache, ri.UserID, r.Reminder{RepeatMode: ri.Text})
	mr.Text = "Value?"
	mr.Keyboard = common.CompileCancelKeyboard()
	return mr
}

func ReadReminderValue(ri *r.ReminderInput) (mr models.ModuleResponse) {
	value64, err := strconv.ParseUint(ri.Text, 10, 8)
	if err != nil {
		mr.Text = "Bad value"
		mr.Error = true
		mr.Keyboard = common.CompileCancelKeyboard()
		return mr
	}
	value8 := uint8(value64)
	AppendCache(ri.ReminderCache, ri.UserID, r.Reminder{RepeatValue: value8})
	mr.Text = "Done"
	mr.EndChat = true
	mr.Keyboard = common.CompileDefaultKeyboard()
	WriteReminder(ri.ReminderCache, ri.UserID)
	DeleteReminderCache(ri.ReminderCache, ri.UserID)
	return mr
}

func GetReminders(userID int64) (mr models.ModuleResponse) {
	reminders := LoadReminders()
	for _, reminder := range reminders {
		if reminder.ReminderData.UserID == userID {
			mr.Text += FormatReminder(reminder.ReminderData)
		}
	}
	if len(mr.Text) == 0 {
		mr.Text = "No reminders found."
	}
	return mr
}

func DeleteReminderQuery(userID int64) (mr models.ModuleResponse) {
	var counter int
	response := "Send me the number of reminder to delete:\n"
	mr.Text += response
	reminders := LoadReminders()
	for _, reminder := range reminders {
		if reminder.ReminderData.UserID == userID {
			counter += 1
			mr.Text += fmt.Sprintf("(%v) ", counter)
			mr.Text += FormatReminder(reminder.ReminderData)
		}
	}
	if mr.Text == response {
		mr.Text = "No reminders found."
		mr.Error = true
		mr.Keyboard = common.CompileDefaultKeyboard()
	} else {
		mr.Keyboard = common.CompileCancelKeyboard()
	}
	return mr
}

func DeleteReminderConfirm(input string, userID int64) (mr models.ModuleResponse) {
	number, err := strconv.Atoi(input)
	if (err != nil) || (number <= 0) {
		mr.Text = "Bad value"
		mr.Error = true
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
		mr.Text = "Bad value"
		mr.Error = true
		mr.Keyboard = common.CompileCancelKeyboard()
		return mr
	}
	if !DeleteReminder(userReminders[number].FileName) {
		mr.Text = "Failed to delete the file"
		mr.Error = true
		mr.Keyboard = common.CompileCancelKeyboard()
		return mr
	}
	log.Printf("Reminder deleted - %v.", userReminders[number].ReminderData)
	mr.Text = "Done"
	mr.EndChat = true
	mr.Keyboard = common.CompileDefaultKeyboard()
	return mr
}
