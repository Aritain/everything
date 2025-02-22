package common

import (
    "net/http"
    "encoding/json"
    "encoding/xml"
)

func GetRequest[T any](url string, mode string, params map[string]string, headers map[string]string) (T, bool) {
    var results T

    client := &http.Client{}
    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        return results, true
    }

    reqParams := req.URL.Query()
    for k, v := range params {
        reqParams.Add(k, v)
    }
    req.URL.RawQuery = reqParams.Encode()

    for k, v := range headers {
        req.Header.Add(k, v)
    }

    resp, err := client.Do(req)
    if err != nil {
        return results, true
    }

    if mode == "xml" {
        err = xml.NewDecoder(resp.Body).Decode(&results)
    } else if mode == "json" {
        err = json.NewDecoder(resp.Body).Decode(&results)
    }
    if err != nil {
        return results, true
    }

    return results, false
}
