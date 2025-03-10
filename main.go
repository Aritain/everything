package main

import (
	"log"
	"slices"
	"strings"

	"everything/common"
	c "everything/config"
	"everything/models"
	r "everything/models/reminder"
	"everything/reminder"
	"everything/tfl"
	"everything/weather"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	config, err := c.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load the config.")
	}

	bot, err := tgbotapi.NewBotAPI(config.TGToken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	var userChats []models.SavedChat
	var reminderCache []r.Reminder
	var userID int64
	var chatPath string
	var chatStage int8
	var text string
	remindCreatePath := "create_reminder"
	remindDeletePath := "delete_reminder"
	// Create chan for telegram updates
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates := bot.GetUpdatesChan(ucfg)

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
		usedKeyboard := common.CompileDefaultKeyboard()

		msg := tgbotapi.NewMessage(userID, "")
		// Cancel ongoing conversation and purge cache
		if text == "Cancel" {
			common.EndChat(&userChats, userID)
			reminder.DeleteReminderCache(&reminderCache, userID)
			mr.ResponseText = "Ok"
		}
		chatPath, chatStage = common.FetchUser(&userChats, userID)

		// /create_reminder path
		if chatPath == remindCreatePath {
			usedKeyboard = common.CompileCancelKeyboard()
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
				usedKeyboard = common.CompileYesNoKeyboard()
			case 2:
				mr = reminder.ReadReminderRepeat(&reminderInput)
				if mr.ResponseCode {
					common.EndChat(&userChats, userID)
					usedKeyboard = common.CompileDefaultKeyboard()
				} else {
					usedKeyboard = common.CompileReminderModeKeyboard()
				}
			case 3:
				mr = reminder.ReadReminderMode(&reminderInput)
			case 4:
				mr = reminder.ReadReminderValue(&reminderInput)
				if !mr.ResponseCode {
					common.EndChat(&userChats, userID)
					usedKeyboard = common.CompileDefaultKeyboard()
				}
			}
			if !mr.ResponseCode {
				common.IncrementStage(&userChats, userID)
			}
		}
		// /delete_reminder path
		if chatPath == remindDeletePath {
			mr = reminder.DeleteReminderConfirm(text, userID)
			usedKeyboard = common.CompileCancelKeyboard()
			if !mr.ResponseCode {
				common.EndChat(&userChats, userID)
				usedKeyboard = common.CompileDefaultKeyboard()
			}
		}
		switch text {
		case "tfl":
			mr = tfl.FetchStatus(&config)
		case "weather":
			mr = weather.FetchStatus(&config)
		case remindCreatePath:
			userChats = append(userChats, models.SavedChat{UserID: userID, ChatPath: remindCreatePath, ChatStage: 0})
			reminderCache = append(reminderCache, r.Reminder{UserID: userID})
			mr.ResponseText = "Reminder name?"
			usedKeyboard = common.CompileCancelKeyboard()
		case remindDeletePath:
			userChats = append(userChats, models.SavedChat{UserID: userID, ChatPath: remindDeletePath, ChatStage: 0})
			usedKeyboard = common.CompileCancelKeyboard()
			mr = reminder.DeleteReminderQuery(userID)
		case "get_reminders":
			mr = reminder.GetReminders(userID)
		case "help":
			mr.ResponseText = "Get help" // Make this into a function in common
		}

		// If user message did not match with anything
		if len(mr.ResponseText) == 0 {
			mr.ResponseText = "Get help"
		}
		msg.Text = mr.ResponseText
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = usedKeyboard
		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
