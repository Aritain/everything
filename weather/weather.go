package weather

import (
    "fmt"

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
    //APIResponse, mr.ResponseCode = GetData()
    APIResponse, mr.ResponseCode = common.GetRequest[w.WeatherStatus](
        config.WeatherEndpoint,
        "json",
        params,
        headers,
    )
    fmt.Println(APIResponse)
    mr.ResponseText, mr.ResponseCode = "123", false
    return mr
}
