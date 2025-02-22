package tfl

import (
    "fmt"

    "everything/models"
    "everything/common"
    tfl "everything/models/tfl"
)


func FetchStatus(config *models.Config) (mr models.ModuleResponse) {
    var APIResponse tfl.ArrayOfLineStatus
    var responceData []tfl.TFLParsed
    params := map[string]string{}
    headers := map[string]string{
        "app_key": config.TFLToken,
        "User-Agent": config.TFLAgent,
    }
    trackedLines := []string{"🟪 Elizabeth Line", "🟩 District", "🟦 Piccadilly", "🟥 Central"}

    APIResponse, mr.ResponseCode = common.GetRequest[tfl.ArrayOfLineStatus](
        config.TFLEndpoint,
        "xml",
        params, headers,
    )
    if mr.ResponseCode {
        return mr
    }

    for _, line := range trackedLines {
        for _, entry := range APIResponse.Lines {
            // 5: for skipping color square
            if line[5:] == entry.Line.Name {
                responceData = append(
                    responceData,
                    tfl.TFLParsed{Line: line, Status: entry.Status.Description,
                })
            }
        }
    }

    for _, elem := range responceData {
        mr.ResponseText += fmt.Sprintf("%s - *%s*\n", elem.Line, elem.Status)
    }
    mr.ResponseText += "https://tfl.gov.uk/"
    return mr
}
