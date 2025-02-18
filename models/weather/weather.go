package models

type WeatherStatus struct {
    Forecast Forecast `json:"forecast"`
}

type Forecast struct {
    Forecasts []DayForecast `json:"forecastday"`
}

type DayForecast struct {
    Today Today `json:"day"`
}

type Today struct {
    MaxT       float32 `json:"maxtemp_c"`
    MinT       float32 `json:"mintemp_c"`
    MaxWind    float32 `json:"maxwind_kph"`
    Precip     float32 `json:"totalprecip_mm"`
    RainChance int     `json:"daily_will_it_rain"`
}
