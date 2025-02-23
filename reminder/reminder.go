package reminder 

import (
    time

    r "everything/models/reminder"
)

func StartReminder(reminderCache *[]r.Reminder, text string, userID int64) (mr models.ModuleResponse) {
    reminderCache = append(reminderCache, r.Reminder{UserID: ChatID, ReminderText: text})
    mr.ResponseText = "When?"
    return mr
}

func ReadReminderName(reminderCache *[]r.Reminder, text string, userID int64) (mr models.ModuleResponse) {
    // verify time
    return mr
}