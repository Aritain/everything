package main

import (
    "log"
    "slices"

    "everything/models"
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
    remindCreatePath := "create_reminder"
    remindDeletePath := "delete_reminder"
    // Create chan for telegram updates
    var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
    ucfg.Timeout = 60
    updates := bot.GetUpdatesChan(ucfg)

    for update := range updates {
        var mr models.ModuleResponse
        userID = update.Message.Chat.ID
        chatPath, chatStage = common.FetchUser(&userChats, ChatID)
        if !slices.Contains(config.BotAdmins, userID) {
            continue
        }

        if update.Message == nil { // ignore any non-Message updates
            continue
        }

        msg := tgbotapi.NewMessage(userID, "")
        // TODO - implement cancel
        // /create_reminder path
        if chatPath == remindCreatePath {
            common.EndChat(&userChats, ChatID)
            switch chatStage {
                case 0:
                    mr = reminder.StartReminder(&reminderCache, update.Message.Text ,userID)
                    common.IncrementStage(&userChats, userID)
                case 1:
                    mr = reminder.ProcessTime(&reminderCache, update.Message.Text ,userID)
                    if !mr.ResponseCode {
                        common.IncrementStage(&userChats, userID)
                    }
            }
        }
        // /delete_reminder path
        if chatPath == remindDeletePath {
            common.EndChat(&userChats, ChatID)
            switch chatStage {
            }
        }
        switch update.Message.Command() {
            // TODO - implement help
            case "tfl":
                mr = tfl.FetchStatus(&config)
            case "weather":
                mr = weather.FetchStatus(&config)
            case "create_reminder":
                userChats = append(userChats, models.SavedChat{userID, remindCreatePath, 0})
                mr.ResponseText = "Reminder name?"
            case "delete_reminder"
                userChats = append(userChats, models.SavedChat{userID, remindDeletePath, 0})
                mr = reminder.DeleteReminder(userID)
        }

        if mr.ResponseCode {
            mr.ResponseText = "Failed to process the request."
        }
        msg.Text = mr.ResponseText
        msg.ParseMode = "Markdown"
        if _, err := bot.Send(msg); err != nil {
            log.Panic(err)
        }
    }
}
