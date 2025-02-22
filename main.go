package main

import (
    "log"
    "slices"

    "everything/models"
    "everything/tfl"
    "everything/weather"
    c "everything/config"

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

    var mr models.ModuleResponse
    // Create chan for telegram updates
    var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
    ucfg.Timeout = 60
    updates := bot.GetUpdatesChan(ucfg)

    for update := range updates {
        if !slices.Contains(config.BotAdmins, update.Message.Chat.ID) {
            continue
        }

        if update.Message == nil { // ignore any non-Message updates
            continue
        }

        if !update.Message.IsCommand() { // ignore any non-command Messages
            continue
        }

        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
        switch {
            case update.Message.Command() == "tfl":
                mr = tfl.FetchStatus(&config)
            case update.Message.Command() == "weather":
                mr = weather.FetchStatus(&config)
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
