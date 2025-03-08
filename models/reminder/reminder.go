package models

import (
    "time"
)

type Reminder struct {
    UserID       int64      `json:"UserID"`
    ReminderText string     `json:"ReminderText"`
    NextReminder time.Time  `json:"NextReminder"`
    RepeatToggle bool       `json:"RepeatToggle"`
    RepeatMode   string     `json:"RepeatMode"` // day, week, month, year
    RepeatValue  uint8      `json:"RepeatValue"`
}

type ReminderInput struct {
    ReminderCache *[]Reminder
    Text          string
    UserID        int64
}
