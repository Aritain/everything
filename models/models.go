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
	TGToken          string  `mapstructure:"TG_TOKEN"`
	BotAdmins        []int64 `mapstructure:"BOT_ADMINS"`
	BotDebug         bool    `mapstructure:"BOT_DEBUG"`
	TFLToken         string  `mapstructure:"TFL_TOKEN"`
	TFLEndpoint      string  `mapstructure:"TFL_ENDPOINT"`
	TFLAgent         string  `mapstructure:"TFL_AGENT"` // Call fails without this
	WeatherToken     string  `mapstructure:"WEATHER_TOKEN"`
	WeatherLocation  string  `mapstructure:"WEATHER_LOCATION"`
	WeatherEndpoint  string  `mapstructure:"WEATHER_ENDPOINT"`
	ReminderDir      string  `mapstructure:"REMINDER_DIR"`
	TimezoneLocation string  `mapstructure:"TIMZONE_LOCATION"`
	CodesEndpoint    string  `mapstructure:"CODES_ENDPOINT"`
	CodesDir         string  `mapstructure:"CODES_DIR"`
	CodesURL         string  `mapstructure:"CODES_URL"`
	GoogleToken      string  `mapstructure:"GOOGLE_TOKEN"`
	GoogleDirId      string  `mapstructure:"GOOGLE_DIR_ID"`
}

type SavedChat struct {
	UserID    int64
	ChatPath  string
	ChatStage int8
}

type TGMessage struct {
	TGToken   string
	UserID    int64
	Text      string
	ParseMode string
	Keyboard  t.InlineKeyboardMarkup
}
