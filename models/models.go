package models

type ModuleResponse struct {
    ResponseText string
    ResponseCode bool
}

type Config struct {
    TGToken         string  `mapstructure:"TG_TOKEN"`
    BotAdmins       []int64 `mapstructure:"BOT_ADMINS"`
    TFLToken        string  `mapstructure:"TFL_TOKEN"`
    TFLEndpoint     string  `mapstructure:"TFL_ENDPOINT"`
    TFLAgent        string  `mapstructure:"TFL_AGENT"` // Call fails without this
    WeatherToken    string  `mapstructure:"WEATHER_TOKEN"`
    WeatherEndpoint string  `mapstructure:"WEATHER_ENDPOINT"`
}
