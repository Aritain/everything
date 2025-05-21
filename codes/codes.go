package codes

import (
	"encoding/json"
	"everything/common"
	cfg "everything/config"
	"everything/models"
	c "everything/models/codes"
	"os"
	"slices"
)

const TIMEOUT = 1

func GetCodesUsers() (subscribers c.Subscribers, err error) {
	config := cfg.Get().Config()
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

func SubscribeUser(text string, userID int64) (mr models.ModuleResponse) {
	config := cfg.Get().Config()
	var currentSubscribers c.Subscribers
	var newSubscribers c.Subscribers
	filepath := config.CodesDir + "/" + "subscribers.json"
	currentSubscribers, err := GetCodesUsers()
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
