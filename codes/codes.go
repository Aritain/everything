package codes

import (
	"encoding/json"
	"everything/common"
	"everything/models"
	c "everything/models/codes"
	"fmt"
	"log"
	"os"
	"reflect"
	"slices"
	"strings"
	"time"
)

const TIMEOUT = 1

func GetCodesUsers(config *models.Config) (subscribers c.Subscribers, err error) {
	filepath := config.CodesDir + "/" + "subscribers.json"
	data, err := os.ReadFile(filepath)
	if err == nil {
		_ = json.Unmarshal(data, &subscribers)
	}
	return
}

func AskID(userID int64) (mr models.ModuleResponse) {
	mr.Text = "Provide your ID"
	mr.Keyboard = common.CompileCancelKeyboard()
	return
}

func SubscribeUser(text string, userID int64, config *models.Config) (mr models.ModuleResponse) {
	var currentSubscribers c.Subscribers
	var newSubscribers c.Subscribers
	filepath := config.CodesDir + "/" + "subscribers.json"
	currentSubscribers, err := GetCodesUsers(config)
	if err == nil {
		os.Remove(filepath)
	}
	// Compile current users so we can check slice contains
	var userIDs []int64
	for _, subscriber := range currentSubscribers.Subscribers {
		userIDs = append(userIDs, subscriber.TGID)
	}

	if !slices.Contains(userIDs, userID) {
		currentSubscribers.Subscribers = append(currentSubscribers.Subscribers, c.Subscriber{TGID: userID, UserID: text})
		newSubscribers = currentSubscribers
		mr.Text = "Subscribed to codes."
	} else {
		for _, subscriber := range currentSubscribers.Subscribers {
			if subscriber.TGID != userID {
				newSubscribers.Subscribers = append(newSubscribers.Subscribers, subscriber)
			}
		}
		mr.Text = "Unsubscribed from codes."
	}
	mr.EndChat = true
	file, _ := os.Create(filepath)
	defer file.Close()
	json.NewEncoder(file).Encode(newSubscribers)
	return
}

func FetchCodes(config *models.Config) {
	filepath := config.CodesDir + "/" + "codes.json"
	params := map[string]string{}
	headers := map[string]string{}
	for {
		var CodesResponse c.CodeData
		var CodesStored c.CodeData
		var fetchError bool
		var newCodes []string
		CodesResponse, fetchError = common.GetRequest[c.CodeData](
			config.CodesEndpoint,
			"json",
			params, headers,
		)
		if fetchError {
			time.Sleep(TIMEOUT * time.Hour)
			continue
		}
		data, err := os.ReadFile(filepath)
		if err == nil {
			_ = json.Unmarshal(data, &CodesStored)
		}
		if reflect.DeepEqual(CodesResponse.Codes, CodesStored.Codes) {
			time.Sleep(TIMEOUT * time.Hour)
			continue
		}
		for _, code := range CodesResponse.Codes {
			if !slices.Contains(CodesStored.Codes, code) {
				newCodes = append(newCodes, code.Code)
			}
		}
		users, err := GetCodesUsers(config)
		if err != nil {
			log.Println("No subscribers found, skipping")
			time.Sleep(TIMEOUT * time.Hour)
			continue
		}
		if len(newCodes) != 0 {
			for _, user := range users.Subscribers {
				message := FormatCodes(user.UserID, newCodes, config.CodesURL)
				var tgm models.TGMessage
				tgm.TGToken = config.TGToken
				tgm.UserID = user.TGID
				tgm.Text = message
				tgm.ParseMode = "HTML"
				go common.SendTGMessage(tgm)
			}
		}
		os.Remove(filepath)
		file, _ := os.Create(filepath)
		json.NewEncoder(file).Encode(CodesResponse)
		file.Close()
		time.Sleep(TIMEOUT * time.Hour)
	}
}

func FormatCodes(userID string, codes []string, CodesURL string) (codesFormatted string) {
	for _, code := range codes {
		fmtURL := CodesURL
		fmtURL = strings.Replace(fmtURL, "NEW_CODE", code, -1)
		fmtURL = strings.Replace(fmtURL, "USER_ID", userID, -1)
		codesFormatted += fmt.Sprintf("<a href='%s'>%s</a>\n", fmtURL, code)
	}
	return
}
