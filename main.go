package main

import (
	"log"
	"slices"
	"strings"

	"everything/codes"
	"everything/common"
	cfg "everything/config"
	"everything/entry"
	"everything/models"
	e "everything/models/entry"
	n "everything/models/notes"
	r "everything/models/reminder"
	"everything/notes"
	"everything/reminder"
	"everything/tfl"
	"everything/weather"

	t "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	if err := cfg.Initialize(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	config := cfg.Get().Config()

	bot, err := t.NewBotAPI(config.TGToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = config.BotDebug

	var userChats []models.SavedChat
	var reminderCache []r.Reminder
	var noteCache []n.FileSelector
	var entryCache []e.Entry

	remindCreatePath := "create_reminder"
	remindDeletePath := "delete_reminder"
	entryPath := "create_entry"
	notesPath := "notes"
	codesPath := "codes_subscribe"
	// Create chan for telegram updates
	var ucfg t.UpdateConfig = t.NewUpdate(0)
	ucfg.Timeout = 60
	updates := bot.GetUpdatesChan(ucfg)
	go reminder.WatchReminders()
	go codes.FetchCodes()

	for update := range updates {
		var userID int64
		var chatPath string
		var chatStage int8
		var text string
		var mr models.ModuleResponse
		var tgm models.TGMessage
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
				mr = reminder.ReadReminderTime(&reminderInput)
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
		// /codes_subscribe path
		if chatPath == codesPath {
			mr = codes.SubscribeUser(text, userID)
		}
		// /notes path
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
				mr = tfl.FetchStatus()
			case "weather":
				mr = weather.FetchStatus()
			case codesPath:
				userChats = append(userChats, models.SavedChat{UserID: userID, ChatPath: codesPath, ChatStage: 0})
				mr = codes.AskID(userID)
			case "get_reminders":
				mr = reminder.GetReminders(userID)
			case notesPath:
				userChats = append(userChats, models.SavedChat{UserID: userID, ChatPath: notesPath, ChatStage: 0})
				mr = notes.ListFiles()
			case entryPath:
				userChats = append(userChats, models.SavedChat{UserID: userID, ChatPath: entryPath, ChatStage: 0})
				mr = entry.EntryCreationStart(userID, &entryCache)
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
		// If user message did not match with anything -> assume it's reminder path
		if mr.Text == "" {
			reminderCache = append(reminderCache, r.Reminder{UserID: userID})
			userChats = append(userChats, models.SavedChat{UserID: userID, ChatPath: remindCreatePath, ChatStage: 1})
			mr = reminder.ReadReminderName(&r.ReminderInput{ReminderCache: &reminderCache, Text: text, UserID: userID})
		}
		if len(mr.Keyboard.InlineKeyboard) == 0 {
			mr.Keyboard = common.CompileDefaultKeyboard()
		}
		tgm.TGToken = config.TGToken
		tgm.UserID = userID
		tgm.Text = mr.Text
		tgm.ParseMode = "Markdown"
		tgm.Keyboard = mr.Keyboard
		go common.SendTGMessage(tgm)
		if err != nil {
			log.Panic(err)
		}
	}
}
