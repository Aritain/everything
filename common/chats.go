package common

import (
	"everything/models"
)

func EndChat(userChats *[]models.SavedChat, userID int64) {
	for index, elem := range *userChats {
		if elem.UserID == userID {
			*userChats = append((*userChats)[:index], (*userChats)[index+1:]...)
		}
	}
}

func FetchUser(userChats *[]models.SavedChat, userID int64) (string, int8) {
	for _, elem := range *userChats {
		if elem.UserID == userID {
			return elem.ChatPath, elem.ChatStage
		}
	}
	return "", 0
}

func IncrementStage(userChats *[]models.SavedChat, userID int64) {
	for i, v := range *userChats {
		if v.UserID == userID {
			(*userChats)[i].ChatStage += 1
			break
		}
	}
}
