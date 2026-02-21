package tfl

import (
	"fmt"

	"everything/common"
	cfg "everything/config"
	"everything/models"
	tfl "everything/models/tfl"
)

func FetchStatus() (mr models.ModuleResponse) {
	config := cfg.Get().Config()
	var APIResponse tfl.ArrayOfLineStatus
	var responseData []tfl.TFLParsed
	params := map[string]string{}
	headers := map[string]string{
		"app_key":    config.TFLToken,
		"User-Agent": config.TFLAgent,
	}
	trackedLines := []string{"💜 Elizabeth Line", "💚 District", "🩶 Jubilee", "🩵 DLR"}
	APIResponse, mr.Error = common.GetRequest[tfl.ArrayOfLineStatus](
		config.TFLEndpoint,
		"xml",
		params, headers,
	)
	if mr.Error {
		mr.Text = "Failed to fetch TFL data"
		return mr
	}

	for _, line := range trackedLines {
		for _, entry := range APIResponse.Lines {
			// 5: for skipping color square
			if line[5:] == entry.Line.Name {
				responseData = append(
					responseData,
					tfl.TFLParsed{Line: line, Status: entry.Status.Description})
			}
		}
	}

	for _, elem := range responseData {
		mr.Text += fmt.Sprintf("%s - *%s*\n", elem.Line, elem.Status)
	}
	mr.Text += "https://tfl.gov.uk/"
	return mr
}
