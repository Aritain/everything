package tfl

import (
    "log"
    "io/ioutil"
    "os"
    "net/http"
    "encoding/xml"
    "everything/models"
)

const TFL_URL = "https://api.tfl.gov.uk/trackernet/LineStatus"
// API auth fails without User-Agent set
const AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"


func FetchStatus() (tr models.ModuleResponse) {
    var APIResponse ArrayOfLineStatus
    APIResponse, tr.ResponseCode = GetData()
    if tr.ResponseCode {
        return tr
    }
    log.Println(APIResponse)

    tr.ResponseText = "OK"
    return tr
}


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
