package reminder

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	cfg "everything/config"
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
	config := cfg.Get().Config()
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
	log.Printf("Reminder created - %v.", reminder)
}

func LoadReminders() (reminders []r.ReminderFile) {
	config := cfg.Get().Config()
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
			" and every %v %s",
			reminder.RepeatValue,
			reminder.RepeatMode,
		)
		if reminder.RepeatValue > 1 {
			repeatable += "s"
		}
	}
	repeatable += "\n"
	fmtReminder = fmt.Sprintf(
		"Reminder *%s*\n%s%s\n",
		reminder.ReminderText,
		reminder.NextReminder.Format("2006-01-02 15:04"),
		repeatable,
	)
	return fmtReminder
}

func DeleteReminder(filename string) (status bool) {
	config := cfg.Get().Config()
	dir := config.ReminderDir
	err := os.Remove(dir + "/" + filename)
	return err == nil
}

func UpdateReminder(reminder r.Reminder) (newTime time.Time) {
	switch {
	case reminder.RepeatMode == "year":
		newTime = reminder.NextReminder.AddDate(int(reminder.RepeatValue), 0, 0)
	case reminder.RepeatMode == "month":
		newTime = reminder.NextReminder.AddDate(0, int(reminder.RepeatValue), 0)
	case reminder.RepeatMode == "week":
		newTime = reminder.NextReminder.AddDate(0, 0, (7 * int(reminder.RepeatValue)))
	case reminder.RepeatMode == "day":
		newTime = reminder.NextReminder.AddDate(0, 0, int(reminder.RepeatValue))
	}
	return newTime
}

/*
Pure AI code, parses user input as time.Time, following formats are supported
YYYY-MM-DD hh:mm
YYYY-MM-DD hh
MM-DD hh:mm
MM-DD hh
DD hh:mm
DD hh
hh:mm
hh
By default set current year/day/etc if not supplied by the user
If upon using current year/day/etc the date happens to be in the past -
Use next year/day/etc
E.g. If today is 2025-01-15 15:00 and user provides an input of "13", then
The value would be 2025-01-16 13:00
*/
func ParseTime(input string) (time.Time, error) {
	config := cfg.Get().Config()
	// Use a fixed reference to UTC and then convert to desired timezone
	location, err := time.LoadLocation(config.TimezoneLocation)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timezone: %v", err)
	}

	// Get current time in UTC first, then convert to location
	nowUTC := time.Now().UTC()
	now := nowUTC.In(location)

	input = strings.TrimSpace(input)

	// Handle relative formats like 15d, 3h, 2m, 1y
	if len(input) > 1 {
		suffix := input[len(input)-1]
		valueStr := input[:len(input)-1]
		if value, err := strconv.Atoi(valueStr); err == nil {
			switch suffix {
			case 'h':
				return now.Add(time.Duration(value) * time.Hour), nil
			case 'd':
				return now.AddDate(0, 0, value), nil
			case 'm':
				return now.AddDate(0, value, 0), nil
			case 'y':
				return now.AddDate(value, 0, 0), nil
			}
		}
	}

	switch {
	case strings.Contains(input, "-") && strings.Contains(input, ":"):
		// Format: YYYY-MM-DD hh:mm or MM-DD hh:mm
		parts := strings.Split(input, " ")
		if len(parts) != 2 {
			return time.Time{}, fmt.Errorf("invalid format: expected date and time separated by space")
		}

		datePart := parts[0]
		timePart := parts[1]

		if strings.Count(datePart, "-") == 2 {
			// Format: YYYY-MM-DD hh:mm
			t, err := time.ParseInLocation("2006-01-02 15:04", input, location)
			if err != nil {
				return time.Time{}, err
			}
			return t, nil
		} else {
			// Format: MM-DD hh:mm
			datePart = fmt.Sprintf("%d-%s", now.Year(), datePart)
			t, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", datePart, timePart), location)
			if err != nil {
				return time.Time{}, err
			}

			if t.Before(now) {
				t = t.AddDate(1, 0, 0)
			}
			return t, nil
		}

	case strings.Contains(input, "-") && !strings.Contains(input, ":"):
		// Format: YYYY-MM-DD hh or MM-DD hh
		parts := strings.Split(input, " ")
		if len(parts) != 2 {
			return time.Time{}, fmt.Errorf("invalid format: expected date and time separated by space")
		}

		datePart := parts[0]
		hourPart := parts[1]

		if strings.Count(datePart, "-") == 2 {
			t, err := time.ParseInLocation("2006-01-02 15", fmt.Sprintf("%s %s", datePart, hourPart), location)
			if err != nil {
				return time.Time{}, err
			}
			return t, nil
		} else {
			datePart = fmt.Sprintf("%d-%s", now.Year(), datePart)
			t, err := time.ParseInLocation("2006-01-02 15", fmt.Sprintf("%s %s", datePart, hourPart), location)
			if err != nil {
				return time.Time{}, err
			}

			if t.Before(now) {
				t = t.AddDate(1, 0, 0)
			}
			return t, nil
		}

	case strings.Contains(input, " "):
		// Format: DD hh:mm or DD hh
		parts := strings.Split(input, " ")
		if len(parts) != 2 {
			return time.Time{}, fmt.Errorf("invalid format: expected day and time separated by space")
		}

		dayPart := parts[0]
		timePart := parts[1]

		if strings.Contains(timePart, ":") {
			datePart := fmt.Sprintf("%d-%02d-%s", now.Year(), now.Month(), dayPart)
			t, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", datePart, timePart), location)
			if err != nil {
				return time.Time{}, err
			}

			if t.Before(now) {
				t = t.AddDate(0, 1, 0)
			}
			return t, nil
		} else {
			datePart := fmt.Sprintf("%d-%02d-%s", now.Year(), now.Month(), dayPart)
			t, err := time.ParseInLocation("2006-01-02 15", fmt.Sprintf("%s %s", datePart, timePart), location)
			if err != nil {
				return time.Time{}, err
			}

			if t.Before(now) {
				t = t.AddDate(0, 1, 0)
			}
			return t, nil
		}

	case strings.Contains(input, ":"):
		// Format: hh:mm
		parts := strings.Split(input, ":")
		if len(parts) != 2 {
			return time.Time{}, fmt.Errorf("invalid time format")
		}

		hour, err := strconv.Atoi(parts[0])
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid hour: %v", err)
		}
		minute, err := strconv.Atoi(parts[1])
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid minute: %v", err)
		}

		// Create time using today's date in the given location
		// Directly construct the time with explicit fields to avoid any DST issues
		year, month, day := now.Date()
		t := time.Date(year, month, day, hour, minute, 0, 0, location)

		// Compare with now using Unix timestamps to ensure accurate comparison
		if t.Unix() < now.Unix() {
			t = t.AddDate(0, 0, 1)
		}
		return t, nil

	default:
		// Format: hh
		hour, err := strconv.Atoi(input)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid hour: %v", err)
		}

		// Create time using today's date in the given location
		year, month, day := now.Date()
		t := time.Date(year, month, day, hour, 0, 0, 0, location)

		// Compare with now using Unix timestamps to ensure accurate comparison
		if t.Unix() < now.Unix() {
			t = t.AddDate(0, 0, 1)
		}
		return t, nil
	}
}
