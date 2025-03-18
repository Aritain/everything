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
	filepath := GenerateFilePath(config)
	data, err := os.ReadFile(filepath)
	if err == nil {
		_ = json.Unmarshal(data, &subscribers)
	}
	return
}

func GenerateFilePath(config *models.Config) (filepath string) {
	dir := config.CodesDir
	filename := "subscribers.json"
	filepath = dir + "/" + filename
	return
}

func SubscribeUser(userID int64, config *models.Config) (mr models.ModuleResponse) {
	var currentSubscribers c.Subscribers
	var newSubscribers c.Subscribers
	filepath := GenerateFilePath(config)
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
	dir := config.CodesDir
	filename := "codes.json"
	filepath := dir + "/" + filename
	params := map[string]string{}
	headers := map[string]string{}
	for {
		var CodesResponce c.CodeData
		var CodesStored c.CodeData
		var fetchError bool
		newCodes := "New codes:\n\n"
		CodesResponce, fetchError = common.GetRequest[c.CodeData](
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
		if reflect.DeepEqual(CodesResponce.Codes, CodesStored.Codes) {
			time.Sleep(TIMEOUT * time.Hour)
			continue
		}
		for _, code := range CodesResponce.Codes {
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
		for _, user := range users.Subscriber {
			common.SendTGMessage(user, newCodes)
		}
		os.Remove(filepath)
		file, _ := os.Create(filepath)
		json.NewEncoder(file).Encode(CodesResponce)
		file.Close()
		time.Sleep(TIMEOUT * time.Hour)
	}
}
