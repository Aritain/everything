package reminder

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	c "everything/config"
	r "everything/models/reminder"
)

func DeleteReminderCache(savedReminders *[]r.Reminder, userID int64) {
	for index, elem := range *savedReminders {
		if elem.UserID == userID {
			*savedReminders = append((*savedReminders)[:index], (*savedReminders)[index+1:]...)
		}
	}
}

func AppendCache(savedReminders *[]r.Reminder, userID int64, data r.Reminder) {
	for i, v := range *savedReminders {
		if v.UserID == userID {
			switch {
			case data.ReminderText != "":
				(*savedReminders)[i].ReminderText = data.ReminderText
			case !data.NextReminder.IsZero():
				(*savedReminders)[i].NextReminder = data.NextReminder
			case data.RepeatToggle:
				(*savedReminders)[i].RepeatToggle = data.RepeatToggle
			case data.RepeatMode != "":
				(*savedReminders)[i].RepeatMode = data.RepeatMode
			case data.RepeatValue != 0:
				(*savedReminders)[i].RepeatValue = data.RepeatValue
			}
			break
		}
	}
}

func WriteReminder(savedReminders *[]r.Reminder, userID int64) {
	var reminder r.Reminder
	config, _ := c.LoadConfig()
	dir := config.ReminderDir
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	for _, elem := range *savedReminders {
		if elem.UserID == userID {
			reminder = elem
			break
		}
	}
	filename := fmt.Sprintf("%s%s_%v.json", dir, timestamp, userID)
	file, _ := os.Create(filename)
	defer file.Close()
	json.NewEncoder(file).Encode(reminder)
}

func LoadReminders() (reminders []r.ReminderFile) {
	config, _ := c.LoadConfig()
	dir := config.ReminderDir
	files, _ := os.ReadDir(dir)

	for _, file := range files {
		data, _ := os.ReadFile(dir + "/" + file.Name())
		var reminder r.Reminder
		_ = json.Unmarshal(data, &reminder)
		reminders = append(reminders, r.ReminderFile{FileName: file.Name(), ReminderData: reminder})
	}

	return reminders
}

func FormatReminder(reminder r.Reminder) (fmtReminder string) {
	var repeatable string
	if reminder.RepeatToggle {
		repeatable = fmt.Sprintf(
			" Repeatable each %v %ss.",
			reminder.RepeatValue,
			reminder.RepeatMode,
		)
	}
	repeatable += "\n"
	fmtReminder = fmt.Sprintf(
		"Reminder %s at %s.%s",
		reminder.ReminderText,
		reminder.NextReminder.Format("2006-01-02 15:04"),
		repeatable,
	)
	return fmtReminder
}

/*
Pure AI code, parses user input as time.Time, following formats are supported
YYYY-MM-DD hh:mm
MM-DD hh:mm
DD hh:mm
hh:mm
hh
By default set current year/day/etc if not supplied by the user
If upon using current year/day/etc the date happens to be in the past -
Use next year/day/etc
E.g. If today is 2025-01-15 15:00 and user provides an input of "13", then
The value would be 2025-01-16 13:00
*/
func ParseTime(input string) (time.Time, error) {
	now := time.Now()
	input = strings.TrimSpace(input)
	switch {
	case strings.Contains(input, "-") && strings.Contains(input, ":"):
		parts := strings.Split(input, " ")
		if len(parts) != 2 {
			return time.Time{}, fmt.Errorf("Failed")
		}

		datePart := parts[0]
		timePart := parts[1]

		if strings.Count(datePart, "-") == 2 {
			return time.Parse("2006-01-02 15:04", input)
		} else {
			datePart = fmt.Sprintf("%d-%s", now.Year(), datePart)
			t, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", datePart, timePart))
			if err != nil {
				return time.Time{}, err
			}
			if t.Before(now) {
				t = t.AddDate(1, 0, 0)
			}
			return t, nil
		}

	case strings.Contains(input, " "):
		parts := strings.Split(input, " ")
		if len(parts) != 2 {
			return time.Time{}, fmt.Errorf("Failed")
		}

		dayPart := parts[0]
		timePart := parts[1]

		datePart := fmt.Sprintf("%d-%02d-%s", now.Year(), now.Month(), dayPart)
		t, err := time.Parse("2006-01-02 15:04", fmt.Sprintf("%s %s", datePart, timePart))
		if err != nil {
			return time.Time{}, err
		}
		if t.Before(now) {
			t = t.AddDate(0, 1, 0)
		}
		return t, nil

	case strings.Contains(input, ":"):
		t, err := time.Parse("15:04", input)
		if err != nil {
			return time.Time{}, err
		}
		t = time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
		if t.Before(now) {
			t = t.AddDate(0, 0, 1)
		}
		return t, nil

	default:
		t, err := time.Parse("15", input)
		if err != nil {
			return time.Time{}, err
		}
		t = time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), 0, 0, 0, now.Location())
		if t.Before(now) {
			t = t.AddDate(0, 0, 1)
		}
		return t, nil
	}
}
