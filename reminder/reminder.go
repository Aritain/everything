package reminder 

import (
    "slices"
    "time"

    "everything/models"
    r "everything/models/reminder"
)


func ReadReminderName(ri *r.ReminderInput) (mr models.ModuleResponse) {
    if len(ri.Text) > 30 {
        mr.ResponseText = "Too long"
        mr.ResponseCode = true
        return mr
    }
    ri.ReminderCache = append(ri.ReminderCache, r.Reminder{UserID: ri.UserID, ReminderText: ri.Text})
    mr.ResponseText = "When?"
    return mr
}


func ReadReminderTime(ri *r.ReminderInput) (mr models.ModuleResponse) {
    // verify time
    now := time.Now()
    ri.ReminderCache = append(ri.ReminderCache, r.Reminder{UserID: ri.UserID, NextReminder: now})
    mr.ResponseText = "Repeat?"
    return mr
}


func ReadReminderRepeat(ri *r.ReminderInput) (mr models.ModuleResponse) {
    ri.ReminderCache = append(ri.ReminderCache, r.Reminder{UserID: ri.UserID, RepeatToggle: true})
    mr.ResponseText = "Mode?"
    return mr
}


func ReadReminderMode(ri *r.ReminderInput) (mr models.ModuleResponse) {
    allowedResponces := []string{"day", "week", "month", "year"}
    if !slices.Contains(allowedResponces, ri.Text) {
        mr.ResponseText = "day/week/month/year"
        mr.ResponseCode = true
        return mr
    }
    ri.ReminderCache = append(ri.ReminderCache, r.Reminder{UserID: ri.UserID, RepeatMode: ri.Text})
    mr.ResponseText = "Value?"
    return mr
}


func ReadReminderValue(ri *r.ReminderInput) (mr models.ModuleResponse) {
    ri.ReminderCache = append(ri.ReminderCache, r.Reminder{UserID: ri.UserID, RepeatValue: 1})
    mr.ResponseText = "Done"
    return mr
}
