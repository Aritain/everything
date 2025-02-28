package models

import (
    "time"
)

type Reminder struct {
    UserID       int64
    ReminderText string
    NextReminder time.Time
    RepeatToggle bool
    RepeatMode   string    // day, week, month, year
    RepeatValue  uint8
}

type ReminderInput struct {
    ReminderCache []Reminder
    Text          string
    UserID        int64
}
