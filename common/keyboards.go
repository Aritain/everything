package common

import (
	t "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CompileYesNoKeyboard() t.InlineKeyboardMarkup {
	var keyboard = t.NewInlineKeyboardMarkup(
		t.NewInlineKeyboardRow(
			t.NewInlineKeyboardButtonData("Yes", "Yes"),
			t.NewInlineKeyboardButtonData("No", "No"),
		),
		t.NewInlineKeyboardRow(
			t.NewInlineKeyboardButtonData("Cancel", "Cancel"),
		),
	)
	return keyboard
}

func CompileCancelKeyboard() t.InlineKeyboardMarkup {
	var keyboard = t.NewInlineKeyboardMarkup(
		t.NewInlineKeyboardRow(
			t.NewInlineKeyboardButtonData("Cancel", "Cancel"),
		),
	)
	return keyboard
}

func CompileDefaultKeyboard() t.InlineKeyboardMarkup {
	var keyboard = t.NewInlineKeyboardMarkup(
		t.NewInlineKeyboardRow(
			t.NewInlineKeyboardButtonData("TFL", "/tfl"),
			t.NewInlineKeyboardButtonData("Weather", "/weather"),
		),
		t.NewInlineKeyboardRow(
			t.NewInlineKeyboardButtonData("Create Reminder", "/create_reminder"),
		),
		t.NewInlineKeyboardRow(
			t.NewInlineKeyboardButtonData("Delete Reminder", "/delete_reminder"),
		),
		t.NewInlineKeyboardRow(
			t.NewInlineKeyboardButtonData("Show Reminders", "/get_reminders"),
		),
	)
	return keyboard
}

func CompileReminderModeKeyboard() t.InlineKeyboardMarkup {
	var keyboard = t.NewInlineKeyboardMarkup(
		t.NewInlineKeyboardRow(
			t.NewInlineKeyboardButtonData("Day", "day"),
			t.NewInlineKeyboardButtonData("Week", "week"),
		),
		t.NewInlineKeyboardRow(
			t.NewInlineKeyboardButtonData("Month", "month"),
			t.NewInlineKeyboardButtonData("Year", "year"),
		),
		t.NewInlineKeyboardRow(
			t.NewInlineKeyboardButtonData("Cancel", "Cancel"),
		),
	)
	return keyboard
}
