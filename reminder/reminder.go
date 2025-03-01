package reminder 

import (
    "fmt"
    "slices"
    "strconv"

    "everything/models"
    r "everything/models/reminder"
)

// TODO -
// 1. remove comments
// 2. implement keyboard
func ReadReminderName(ri *r.ReminderInput) (mr models.ModuleResponse) {
    // Done
    if len(ri.Text) > 30 {
        mr.ResponseText = "Too long"
        mr.ResponseCode = true
        return mr
    }
    AppendCache(ri.ReminderCache, ri.UserID, r.Reminder{ReminderText: ri.Text})
    mr.ResponseText = "When?"
    return mr
}


func ReadReminderTime(ri *r.ReminderInput) (mr models.ModuleResponse) {
    // TODO - keyboard here
    reminderTime, err := ParseTime(ri.Text)
    if err != nil {
        mr.ResponseText = "Failed to read time provided"
        mr.ResponseCode = true
        return mr
    }
    AppendCache(ri.ReminderCache, ri.UserID, r.Reminder{NextReminder: reminderTime})
    mr.ResponseText = "Repeat?"
    return mr
}


func ReadReminderRepeat(ri *r.ReminderInput) (mr models.ModuleResponse) {
    // TODO - call WriteReminder
    if ri.Text == "No" {
        mr.ResponseText = "Done"
        mr.ResponseCode = true
        fmt.Println(ri.ReminderCache)
        DeleteReminderCache(ri.ReminderCache, ri.UserID)
        return mr
    }
    AppendCache(ri.ReminderCache, ri.UserID, r.Reminder{RepeatToggle: true})
    mr.ResponseText = "Mode?"
    return mr
}


func ReadReminderMode(ri *r.ReminderInput) (mr models.ModuleResponse) {
    // TODO - keyboard here
    allowedResponces := []string{"day", "week", "month", "year"}
    if !slices.Contains(allowedResponces, ri.Text) {
        mr.ResponseText = "day/week/month/year"
        mr.ResponseCode = true
        return mr
    }
    AppendCache(ri.ReminderCache, ri.UserID, r.Reminder{RepeatMode: ri.Text})
    mr.ResponseText = "Value?"
    return mr
}


func ReadReminderValue(ri *r.ReminderInput) (mr models.ModuleResponse) {
    // TODO - call WriteReminder
    value64, err := strconv.ParseUint(ri.Text, 10, 8)
    if err != nil {
        mr.ResponseText = "Bad value"
        mr.ResponseCode = true
        return mr
    }
    value8 := uint8(value64)
    AppendCache(ri.ReminderCache, ri.UserID, r.Reminder{RepeatValue: value8})
    mr.ResponseText = "Done"
    fmt.Println(ri.ReminderCache)
    DeleteReminderCache(ri.ReminderCache, ri.UserID)
    return mr
}
