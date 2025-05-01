package main

import (
	"log"
	"slices"
	"strings"

	"everything/codes"
	"everything/common"
	c "everything/config"
	"everything/models"
	n "everything/models/notes"
	r "everything/models/reminder"
	"everything/notes"
	"everything/reminder"
	"everything/tfl"
	"everything/weather"

	t "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	config, err := c.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load the config.")
	}

	bot, err := t.NewBotAPI(config.TGToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = config.BotDebug

	var userChats []models.SavedChat
	var reminderCache []r.Reminder
	var noteCache []n.FileSelector
	var userID int64
	var chatPath string
	var chatStage int8
	var text string
	remindCreatePath := "create_reminder"
	remindDeletePath := "delete_reminder"
	notesPath := "notes"
	codesPath := "codes_subscribe"
	// Create chan for telegram updates
	var ucfg t.UpdateConfig = t.NewUpdate(0)
	ucfg.Timeout = 60
	updates := bot.GetUpdatesChan(ucfg)
	go reminder.WatchReminders(&config)
	go codes.FetchCodes(&config)

	for update := range updates {
		if (update.Message == nil) && (update.CallbackQuery == nil) { // ignore any non-Message updates
			continue
		}
		// Treat CallbackQueries the same as a message from user
		if update.CallbackQuery != nil {
			callback := update.CallbackQuery
			userID = callback.Message.Chat.ID
			text = callback.Data
		} else {
			userID = update.Message.Chat.ID
			text = update.Message.Text
		}
		// Ignore messages from non-whitelisted users
		if !slices.Contains(config.BotAdmins, userID) {
			continue
		}
		// Workaround to treat actual commands the same as responses
		// From Inline keyboard
		if strings.Contains(text, "/") {
			text = text[1:]
		}
		var mr models.ModuleResponse

		msg := t.NewMessage(userID, "")
		// Cancel ongoing conversation and purge cache
		if text == "Cancel" {
			common.EndChat(&userChats, userID)
			reminder.DeleteReminderCache(&reminderCache, userID)
			mr.Text = "Ok"
		}
		chatPath, chatStage = common.FetchUser(&userChats, userID)

		// /create_reminder path
		if chatPath == remindCreatePath {
			reminderInput := r.ReminderInput{
				ReminderCache: &reminderCache,
				Text:          text,
				UserID:        userID,
			}
			switch chatStage {
			case 0:
				mr = reminder.ReadReminderName(&reminderInput)
			case 1:
				mr = reminder.ReadReminderTime(&reminderInput, &config)
			case 2:
				mr = reminder.ReadReminderRepeat(&reminderInput)
			case 3:
				mr = reminder.ReadReminderMode(&reminderInput)
			case 4:
				mr = reminder.ReadReminderValue(&reminderInput)
			}
			if (!mr.Error) && (!mr.EndChat) {
				common.IncrementStage(&userChats, userID)
			}
		}
		// /delete_reminder path
		if chatPath == remindDeletePath {
			mr = reminder.DeleteReminderConfirm(text, userID)
		}
		// /codes path
		if chatPath == codesPath {
			mr = codes.SubscribeUser(text, userID, &config)
		}
		if chatPath == notesPath {
			switch chatStage {
			case 0:
				mr = notes.SelectFile(text, userID, &noteCache)
			case 1:
				mr = notes.UpdateFile(text, userID, &noteCache)
			}
			if (!mr.Error) && (!mr.EndChat) {
				common.IncrementStage(&userChats, userID)
			}
		}
		// Delete any ongoing chat if got EndChat flag
		if mr.EndChat {
			common.EndChat(&userChats, userID)
		}
		// Default commands if user is not in a conversation
		if chatPath == "" {
			switch text {
			case "tfl":
				mr = tfl.FetchStatus(&config)
			case "weather":
				mr = weather.FetchStatus(&config)
			case codesPath:
				userChats = append(userChats, models.SavedChat{UserID: userID, ChatPath: codesPath, ChatStage: 0})
				mr = codes.AskID(userID)
			case "get_reminders":
				mr = reminder.GetReminders(userID)
			case notesPath:
				userChats = append(userChats, models.SavedChat{UserID: userID, ChatPath: notesPath, ChatStage: 0})
				mr = notes.ListFiles()
			case remindCreatePath:
				userChats = append(userChats, models.SavedChat{UserID: userID, ChatPath: remindCreatePath, ChatStage: 0})
				mr = reminder.ReminderCreationStart(userID, &reminderCache)
			case remindDeletePath:
				mr = reminder.DeleteReminderQuery(userID)
				if !mr.Error {
					userChats = append(userChats, models.SavedChat{UserID: userID, ChatPath: remindDeletePath, ChatStage: 0})
				}
			}
		}
		// If user message did not match with anything
		if mr.Text == "" {
			mr.Text = "Try again."
		}
		if len(mr.Keyboard.InlineKeyboard) == 0 {
			mr.Keyboard = common.CompileDefaultKeyboard()
		}
		msg.Text = mr.Text
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = mr.Keyboard
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
