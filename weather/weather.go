package weather

import (
    "fmt"
    "math"

    "everything/models"
    "everything/common"
    w "everything/models/weather"
)


func FetchStatus(config *models.Config) (mr models.ModuleResponse) {
    var APIResponse w.WeatherStatus
    params := map[string]string{
        "key": config.WeatherToken,
        "q": "London",
        "days": "1",
        "aqi": "no",
        "alerts": "no",
    }
    headers := map[string]string{}

    APIResponse, mr.ResponseCode = common.GetRequest[w.WeatherStatus](
        config.WeatherEndpoint,
        "json",
        params,
        headers,
    )

    weatherData := APIResponse.Forecast.Forecasts[0].Today
    maxT := int(math.Round(float64(weatherData.MaxT)))
    minT := int(math.Round(float64(weatherData.MinT)))
    windStr := int(math.Round(float64(weatherData.MaxWind)))
    windType := CheckWind(weatherData.MaxWind)
    rainType := CheckRain(weatherData.Precip)

    mr.ResponseText = fmt.Sprintf("☀️ Maximum temperature - *%d°*\n", maxT)
    mr.ResponseText += fmt.Sprintf("❄️ Minimum temperature - *%d°*\n", minT)
    mr.ResponseText += fmt.Sprintf("💨 %s wind - *%d* km/h\n", windType, windStr)
    mr.ResponseText += fmt.Sprintf("☔️ %s - *%.2f* mm\n", rainType, weatherData.Precip)
    mr.ResponseText += fmt.Sprintf("🔮 Rain chance - *%d*%%", weatherData.RainChance)
    return mr
}


func CheckWind(windSpeed float32) string {
    if windSpeed < 5.0 {
        return "Calm"
    }
    if (windSpeed >= 5.0) && (windSpeed < 20.0) {
        return "Light"
    }
    if (windSpeed >= 20.0) && (windSpeed < 40.0) {
        return "Moderate"
    }
    return "Strong"
}


func CheckRain(precipitation float32) string {
    if precipitation < 0.1 {
        return "No Rain"
    }
    if (precipitation >= 0.1) && (precipitation < 2.4) {
        return "Drizzle"
    }
    if (precipitation >= 2.4) && (precipitation < 7.5) {
        return "Moderate Rain"
    }
    return "Showers"
}
