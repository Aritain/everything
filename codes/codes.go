package codes

import (
	"everything/common"
	"everything/models"
	c "everything/models/codes"
	"fmt"
)

func FetchCodes(config *models.Config) (mr models.ModuleResponse) {
	var CodesResponce c.CodeData
	params := map[string]string{}
	headers := map[string]string{}
	CodesResponce, mr.Error = common.GetRequest[c.CodeData](
		config.CodeEndpoint,
		"json",
		params, headers,
	)
	mr.Text = "Codes:\n"
	for _, code := range CodesResponce.Codes {
		mr.Text += fmt.Sprintf("%s\n", code.Code)
	}
	return mr
}
