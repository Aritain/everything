package weather

import (
    "fmt"

    "everything/models"
    w "everything/models/weather"
)

// TODO - pass as options
const WEATHER_URL_ONE = "https://api.weatherapi.com/v1/forecast.json?key="
const WEATHER_URL_TWO = "&q=London&days=1&aqi=no&alerts=no"

func FetchStatus() (mr models.ModuleResponse) {
    var APIResponse w.WeatherStatus
    APIResponse, mr.ResponseCode = GetData()
    fmt.Println(APIResponse)
    mr.ResponseText, mr.ResponseCode = "123", false
    return mr
}