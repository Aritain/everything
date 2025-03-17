package models

import (
	t "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ModuleResponse struct {
	Text     string
	Error    bool
	EndChat  bool
	Keyboard t.InlineKeyboardMarkup
}

type Config struct {
	TGToken         string  `mapstructure:"TG_TOKEN"`
	BotAdmins       []int64 `mapstructure:"BOT_ADMINS"`
	BotDebug        bool    `mapstructure:"BOT_DEBUG"`
	TFLToken        string  `mapstructure:"TFL_TOKEN"`
	TFLEndpoint     string  `mapstructure:"TFL_ENDPOINT"`
	TFLAgent        string  `mapstructure:"TFL_AGENT"` // Call fails without this
	WeatherToken    string  `mapstructure:"WEATHER_TOKEN"`
	WeatherEndpoint string  `mapstructure:"WEATHER_ENDPOINT"`
	ReminderDir     string  `mapstructure:"REMINDER_DIR"`
	CodeEndpoint    string  `mapstructure:"CODE_ENDPOINT"`
}

type SavedChat struct {
	UserID    int64
	ChatPath  string
	ChatStage int8
}
