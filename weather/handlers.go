package weather

import (
    "log"
    "io/ioutil"
    "os"
    "net/http"
    "encoding/json"

    w "everything/models/weather"
)
// TODO make this call reusable for other modules
func GetData() (APIResponse w.WeatherStatus, processingError bool) {
    APIKey, status := os.LookupEnv("WEATHER_TOKEN")
    if !status {
        log.Printf("WEATHER_TOKEN env is missing.")
        return APIResponse, true
    }

    url := WEATHER_URL_ONE + APIKey + WEATHER_URL_TWO
    log.Println(url)
    client := &http.Client{}
    req , err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Println(err)
        return APIResponse, true
    }

    response , err := client.Do(req)
    if err != nil {
        log.Println(err)
        return APIResponse, true
    }

    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Println(err)
        return APIResponse, true
    }

    err = json.Unmarshal(body, &APIResponse)
    if err != nil {
        log.Println(err)
        return APIResponse, true
    }

    return APIResponse, processingError
}