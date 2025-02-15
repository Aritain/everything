package tfl

import (
    "log"
    "io/ioutil"
    "os"
    "net/http"
    "encoding/xml"
)

func GetData() (APIResponse ArrayOfLineStatus, processingError bool) {
    APIKey, status := os.LookupEnv("TFL_TOKEN")
    if !status {
        log.Printf("TFL_TOKEN env is missing.")
        return APIResponse, true
    }

    client := &http.Client{}
    req , err := http.NewRequest("GET", TFL_URL, nil)
    if err != nil {
        log.Println(err)
        return APIResponse, true
    }

    req.Header.Add("app_key", APIKey)
    req.Header.Add("User-Agent", AGENT)
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

    err = xml.Unmarshal(body, &APIResponse)
    if err != nil {
        log.Println(err)
        return APIResponse, true
    }

    return APIResponse, processingError
}
