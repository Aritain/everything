package models

import "time"

type Entry struct {
	UserID      int64     `json:"UserID"`
	Date        time.Time `json:"Date"` // Year-Month-Day
	Description string    `json:"Description"`
	Characters  []string  `json:"Characters"`
}
