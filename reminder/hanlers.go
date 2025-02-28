package reminder

import (
    //"everything/models"
    r "everything/models/reminder"
)

func DeleteWatchCache(savedReminders *[]r.Reminder, userID int64) {
	for index, elem := range *savedReminders {
		if elem.UserID == userID {
			*savedReminders = append((*savedReminders)[:index], (*savedReminders)[index+1:]...)
		}
	}	
}

// Not sure what is this
/*func AppendCache(userChats *[]models.SavedChat, userID int64) { // update this
    for i, v := range *userChats {
        if v.UserID == userID {
            *userChats[i].ChatStage += 1
            break
        }
    }
    return
}*/
