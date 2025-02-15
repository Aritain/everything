package tfl

import (
    "fmt"
    "everything/models"
)

const TFL_URL = "https://api.tfl.gov.uk/trackernet/LineStatus"
// API auth fails without User-Agent set
const AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"


func FetchStatus() (mr models.ModuleResponse) {
    var APIResponse ArrayOfLineStatus
    var responceData []TFLParsed
    trackedLines := []string{"🟪 Elizabeth Line", "🟩 District", "🟦 Piccadilly", "🟥 Central"}

    APIResponse, mr.ResponseCode = GetData()
    if mr.ResponseCode {
        return mr
    }

    for _, line := range trackedLines {
        for _, entry := range APIResponse.Lines {
            // 5: for skipping color square
            if line[5:] == entry.Line.Name {
                responceData = append(responceData, TFLParsed{Line: line, Status: entry.Status.Description})
            }
        }
    }

    for _, elem := range responceData {
        mr.ResponseText += fmt.Sprintf("%s - *%s*\n", elem.Line, elem.Status)
    }
    mr.ResponseText += "https://tfl.gov.uk/"
    return mr
}
