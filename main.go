package main

import (
    "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "log"
    "os"
    "strconv"
    "everything/models"
    "everything/tfl"
)

func main() {

    tgToken, status := os.LookupEnv("TG_TOKEN")
    if !status {
        log.Printf("TG_TOKEN env is missing.")
        os.Exit(1)
    }

    tgAdmin, status := os.LookupEnv("BOT_ADMIN")
    if !status {
        log.Printf("BOT_ADMIN env is missing.")
        os.Exit(1)
    }
    tgAdminID, _ := strconv.ParseInt(tgAdmin, 10, 64)

    bot, err := tgbotapi.NewBotAPI(tgToken)
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
        if update.Message.Chat.ID != tgAdminID { // ignore non-admin messages
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
                mr = tfl.FetchStatus()
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
