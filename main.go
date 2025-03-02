package main

import (
    "log"
    "slices"

    "everything/common"
    "everything/models"
    "everything/reminder"
    "everything/tfl"
    "everything/weather"
    c "everything/config"
    r "everything/models/reminder"

    "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
    remindCreatePath := "create_reminder"
    remindDeletePath := "delete_reminder"
    // Create chan for telegram updates
    var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
    ucfg.Timeout = 60
    updates := bot.GetUpdatesChan(ucfg)

    for update := range updates {
        var mr models.ModuleResponse
        userID = update.Message.Chat.ID
        text := update.Message.Text
        if !slices.Contains(config.BotAdmins, userID) {
            continue
        }

        if update.Message == nil { // ignore any non-Message updates
            continue
        }

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
            reminderInput := r.ReminderInput{
                ReminderCache: &reminderCache,
                Text         : text,
                UserID       : userID,
            }
            switch chatStage {
                case 0:
                    mr = reminder.ReadReminderName(&reminderInput)
                case 1:
                    mr = reminder.ReadReminderTime(&reminderInput)
                case 2:
                    mr = reminder.ReadReminderRepeat(&reminderInput)
                    if mr.ResponseCode {common.EndChat(&userChats, userID)}
                case 3:
                    mr = reminder.ReadReminderMode(&reminderInput)
                case 4:
                    mr = reminder.ReadReminderValue(&reminderInput)
                    if !mr.ResponseCode {common.EndChat(&userChats, userID)}
            }
            if !mr.ResponseCode {
                common.IncrementStage(&userChats, userID)
            }
        }
        // /delete_reminder path
        /*if chatPath == remindDeletePath {
            common.EndChat(&userChats, ChatID)
            switch chatStage {
            }
        }*/
        switch update.Message.Command() {
            case "tfl":
                mr = tfl.FetchStatus(&config)
            case "weather":
                mr = weather.FetchStatus(&config)
            case remindCreatePath:
                userChats = append(userChats, models.SavedChat{userID, remindCreatePath, 0})
                reminderCache = append(reminderCache, r.Reminder{UserID: userID})
                mr.ResponseText = "Reminder name?"
            case remindDeletePath:
                userChats = append(userChats, models.SavedChat{userID, remindDeletePath, 0})
                //mr = reminder.DeleteReminder(userID)
            case "help":
                mr.ResponseText = "Get help" // Make this into a function in common
        }

        // If user message did not match with anything
        if len(mr.ResponseText) == 0 {
            mr.ResponseText = "Get help"
        }
        msg.Text = mr.ResponseText
        msg.ParseMode = "Markdown"
        if _, err := bot.Send(msg); err != nil {
            log.Panic(err)
        }
    }
}
