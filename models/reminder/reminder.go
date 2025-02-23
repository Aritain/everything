package models

import (
    time
)

type Reminder struct {
    UserID       int64
    ReminderText string
    NextReminder time.time
    RepeatToggle bool
    RepeatMode   string    // days, weeks, months
}
