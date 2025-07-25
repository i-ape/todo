package todo

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//
// 📆 OVERDUE CHECK
//

// IsOverdue returns true if the date is before today.
func IsOverdue(date string) bool {
	due, err := time.Parse("2006-01-02", date)
	return err == nil && time.Now().After(due)
}

// ParseNaturalDate parses strings like:
// "tomorrow", "in 3 days", "2024-05-20", "fri", etc.
// ParseNaturalDate parses natural language dates into YYYY-MM-DD format
func ParseNaturalDate(input string) (string, error) {
	input = strings.ToLower(strings.TrimSpace(input))
	if input == "" {
		return "", fmt.Errorf("date cannot be empty")
	}

	// Use local time zone for consistency
	now := time.Now().Local()

	// Handle abbreviation shortcuts (e.g., "today", "tomorrow")
	if f, ok := abbreviationMap[input]; ok {
		return f(now), nil
	}

	// Handle relative dates (e.g., "in 2 days", "2 weeks", "next month")
	if matches := relativeDateRegex.FindStringSubmatch(input); matches != nil {
		unit := matches[2]
		num := 1 // Default to 1 for "next <unit>"
		if matches[1] != "" {
			var err error
			num, err = strconv.Atoi(matches[1])
			if err != nil {
				return "", fmt.Errorf("invalid number: %s", matches[1])
			}
		}
		switch unit {
		case "d", "day", "days":
			return now.AddDate(0, 0, num).Format("2006-01-02"), nil
		case "w", "week", "weeks":
			return now.AddDate(0, 0, num*7).Format("2006-01-02"), nil
		case "m", "month", "months":
			return now.AddDate(0, num, 0).Format("2006-01-02"), nil
		}
	}

	// Handle natural language dates (e.g., "Dec 31, 2025", "31 Dec 2025")
	if date, err := parseNaturalLanguageDate(input, now); err == nil {
		return date, nil
	}

	// Try specific date formats
	for _, layout := range []string{
		"2006-01-02", // YYYY-MM-DD
		"02-01-2006", // DD-MM-YYYY
		"01-02-2006", // MM-DD-YYYY
	} {
		if t, err := time.ParseInLocation(layout, input, time.Local); err == nil {
			return t.Format("2006-01-02"), nil
		}
	}

	return "", fmt.Errorf("invalid date format: %s. Use YYYY-MM-DD, DD-MM-YYYY, MM-DD-YYYY, 'today', 'in N days', or 'Dec 31, 2025'", input)
}

//
// ⏰ PARSE DATE + TIME + DURATION SYNTAX
//

// ParseDateTimeDuration parses strings like:
// "friday @ 18:00 for 1h" or "tomorrow @ 10:30"
func ParseDateTimeDuration(input string) (date, timeStr, duration string, err error) {
	main := input
	if at := strings.Index(input, "@"); at != -1 {
		main = strings.TrimSpace(input[:at])
		rest := strings.TrimSpace(input[at+1:])

		parts := strings.Split(rest, "for")
		timeStr = strings.TrimSpace(parts[0])
		if len(parts) > 1 {
			duration = strings.TrimSpace(parts[1])
		}
	}

	// Parse main date (e.g. "tomorrow")
	date, err = ParseNaturalDate(main)
	if err != nil {
		return "", "", "", err
	}

	// Validate time format (HH:MM)
	if timeStr != "" {
		if _, err := time.Parse("15:04", timeStr); err != nil {
			return "", "", "", fmt.Errorf("invalid time format: %s", timeStr)
		}
	}

	// Validate duration (e.g. "30m", "1h")
	if duration != "" {
		if _, err := time.ParseDuration(duration); err != nil {
			return "", "", "", fmt.Errorf("invalid duration: %s", duration)
		}
	}

	return date, timeStr, duration, nil
}



// parseNaturalLanguageDate handles formats like "Dec 31, 2025" or "31 Dec 2025"
func parseNaturalLanguageDate(input string, now time.Time) (string, error) {
    parts := strings.Fields(input)
    if len(parts) < 2 || len(parts) > 3 {
        return "", fmt.Errorf("invalid natural language date")
    }

    var day, month, year int
    var err error

    // Try "Dec 31, 2025" or "31 Dec 2025"
    if monthNum, ok := monthNames[parts[0]]; ok && len(parts) == 3 {
        month = monthNum
        day, err = strconv.Atoi(parts[1])
        if err != nil {
            return "", fmt.Errorf("invalid day: %s", parts[1])
        }
        year, err = strconv.Atoi(parts[2])
        if err != nil {
            return "", fmt.Errorf("invalid year: %s", parts[2])
        }
    } else if monthNum, ok := monthNames[parts[1]]; ok && len(parts) >= 2 {
        month = monthNum
        day, err = strconv.Atoi(parts[0])
        if err != nil {
            return "", fmt.Errorf("invalid day: %s", parts[0])
        }
        if len(parts) == 3 {
            year, err = strconv.Atoi(parts[2])
            if err != nil {
                return "", fmt.Errorf("invalid year: %s", parts[2])
            }
        } else {
            year = now.Year()
        }
    } else {
        return "", fmt.Errorf("invalid month name")
    }

    // Validate date
    t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
    if t.Year() != year || t.Month() != time.Month(month) || t.Day() != day {
        return "", fmt.Errorf("invalid date: %s", input)
    }

    

//
// 🔠 SHORTCUT KEYWORDS MAP
//

var abbreviationMap = map[string]func(time.Time) string{
	// 📆 Day shortcuts
	"td": formatToday, "tdy": formatToday, "today": formatToday,
	"tm": inDays(1), "tmmrw": inDays(1), "next": inDays(1),
	"af": inDays(2), "aft": inDays(2),
	"yd": inDays(-1), "yst": inDays(-1),
	"now": formatToday, "soon": inDays(3), "later": inDays(7),
	"someday": func(t time.Time) string { return "" },

	// 📅 Weekly shortcuts
	"nw": inDays(7), "nxtwk": inDays(7),
	"n2w": inDays(14), "n3w": inDays(21),
	"eowk": nextWeekday(time.Friday),

	// 📅 Monthly
	"nm": inMonths(1), "em": endOfMonth,

	// 🗓️ Weekday names (auto pick next)
	"mon": nextWeekday(time.Monday), "tue": nextWeekday(time.Tuesday),
	"wed": nextWeekday(time.Wednesday), "thu": nextWeekday(time.Thursday),
	"fri": nextWeekday(time.Friday), "sat": nextWeekday(time.Saturday),
	"sun":    nextWeekday(time.Sunday),
	"nxtmon": nextWeekday(time.Monday), "nxfri": nextWeekday(time.Friday),

	// ⏳ Misc
	"eod":  formatToday,
	"bom":  startOfMonth,
	"eom":  endOfMonth,
	"eonm": endOfNextMonth,
	"eow":  endOfWeek,
	"som":  startOfMonth,
	"sonm": startOfNextMonth,
	"ew": func(t time.Time) string {
		return t.AddDate(0, 0, 7-int(t.Weekday())).Format("2006-01-02")
	},
}

//
// 🧰 INTERNAL HELPERS
//

func formatToday(t time.Time) string {
	return t.Format("2006-01-02")
}

func inDays(n int) func(time.Time) string {
	return func(t time.Time) string {
		return t.AddDate(0, 0, n).Format("2006-01-02")
	}
}

func inMonths(n int) func(time.Time) string {
	return func(t time.Time) string {
		return t.AddDate(0, n, 0).Format("2006-01-02")
	}
}

func endOfMonth(t time.Time) string {
	firstOfNext := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	return firstOfNext.AddDate(0, 0, -1).Format("2006-01-02")
}

func nextWeekday(wd time.Weekday) func(time.Time) string {
	return func(t time.Time) string {
		offset := (int(wd) - int(t.Weekday()) + 7) % 7
		if offset == 0 {
			offset = 7
		}
		return t.AddDate(0, 0, offset).Format("2006-01-02")
	}
}

func startOfMonth(t time.Time) string {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).Format("2006-01-02")
}

func endOfWeek(t time.Time) string {
	offset := 6 - int(t.Weekday())
	return t.AddDate(0, 0, offset).Format("2006-01-02")
}

func startOfNextMonth(t time.Time) string {
	next := t.AddDate(0, 1, 0)
	return time.Date(next.Year(), next.Month(), 1, 0, 0, 0, 0, t.Location()).Format("2006-01-02")
}

func endOfNextMonth(t time.Time) string {
	next := t.AddDate(0, 2, 0)
	return time.Date(next.Year(), next.Month(), 0, 0, 0, 0, 0, t.Location()).Format("2006-01-02")
}

// parseNaturalDate handles natural language date inputs
func parseNaturalDate(input string) (string, error) {
	now := time.Now()
	input = strings.ToLower(strings.TrimSpace(input))

	switch input {
	case "today":
		return now.Format("2006-01-02"), nil
	case "tomorrow":
		return now.AddDate(0, 0, 1).Format("2006-01-02"), nil
	case "next week":
		return now.AddDate(0, 0, 7).Format("2006-01-02"), nil
	case "next month":
		return now.AddDate(0, 1, 0).Format("2006-01-02"), nil
	case "next year":
		return now.AddDate(1, 0, 0).Format("2006-01-02"), nil
	default:
		// in N days or weeks
		if strings.HasPrefix(input, "in ") {
			parts := strings.Split(input, " ")
			if len(parts) == 3 {
				num, err := strconv.Atoi(parts[1])
				if err != nil {
					return "", fmt.Errorf("invalid number in relative date")
				}
				switch parts[2] {
				case "day", "days":
					return now.AddDate(0, 0, num).Format("2006-01-02"), nil
				case "week", "weeks":
					return now.AddDate(0, 0, num*7).Format("2006-01-02"), nil
				}
			}
		}

		// try DD-MM-YYYY
		t, err := time.Parse("02-01-2006", input)
		if err == nil {
			return t.Format("2006-01-02"), nil
		}
		// try YYYY-MM-DD
		t, err = time.Parse("2006-01-02", input)
		if err == nil {
			return t.Format("2006-01-02"), nil
		}
		return "", fmt.Errorf("invalid date format or unsupported natural keyword")
	}
}

// SetDueDate assigns a due date to a task
func SetDueDate(input string, dueDate string) error {
	tasks, _ := LoadTasks()
	found := false

	parsedDate, err := parseNaturalDate(dueDate)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(input)
	for i, task := range tasks {
		if (err == nil && task.ID == id) || task.Text == input {
			tasks[i].DueDate = parsedDate
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task not found")
	}

	return SaveTasks(tasks)
}
