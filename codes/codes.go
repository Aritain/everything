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
	"time"
)

const TIMEOUT = 6

func GetCodesUsers(config *models.Config) (subscribers c.Subscribers, err error) {
	filepath := GenerateFilePath(config.CodesDir, "subscribers.json")
	data, err := os.ReadFile(filepath)
	if err == nil {
		_ = json.Unmarshal(data, &subscribers)
	}
	return
}

func GenerateFilePath(dir string, filename string) string {
	return dir + "/" + filename
}

func SubscribeUser(userID int64, config *models.Config) (mr models.ModuleResponse) {
	var currentSubscribers c.Subscribers
	var newSubscribers c.Subscribers
	filepath := GenerateFilePath(config.CodesDir, "subscribers.json")
	currentSubscribers, err := GetCodesUsers(config)
	if err == nil {
		os.Remove(filepath)
	}

	if !slices.Contains(currentSubscribers.Subscriber, userID) {
		currentSubscribers.Subscriber = append(currentSubscribers.Subscriber, userID)
		newSubscribers = currentSubscribers
		mr.Text = "Subscribed to codes."
	} else {
		for _, subscriber := range currentSubscribers.Subscriber {
			if subscriber != userID {
				newSubscribers.Subscriber = append(newSubscribers.Subscriber, subscriber)
			}
		}
		mr.Text = "Unsubscribed from codes."
	}
	file, _ := os.Create(filepath)
	defer file.Close()
	json.NewEncoder(file).Encode(newSubscribers)
	return
}

func FetchCodes(config *models.Config) {
	newCodesHeader := "New codes:\n\n"
	filepath := GenerateFilePath(config.CodesDir, "codes.json")
	params := map[string]string{}
	headers := map[string]string{}
	for {
		var CodesResponse c.CodeData
		var CodesStored c.CodeData
		var fetchError bool
		newCodes := newCodesHeader
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
				newCodes += fmt.Sprintf("*%s*\n", code.Code)
			}
		}
		users, err := GetCodesUsers(config)
		if err != nil {
			log.Println("No subscribers found, skipping")
			time.Sleep(TIMEOUT * time.Hour)
			continue
		}
		if newCodes != newCodesHeader {
			for _, user := range users.Subscriber {
				common.SendTGMessage(user, newCodes)
			}
		}
		os.Remove(filepath)
		file, _ := os.Create(filepath)
		json.NewEncoder(file).Encode(CodesResponse)
		file.Close()
		time.Sleep(TIMEOUT * time.Hour)
	}
}
