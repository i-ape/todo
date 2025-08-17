package todo

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// IsOverdue checks if a date is before today.
func IsOverdue(date string) bool {
	due, err := time.ParseInLocation("2006-01-02", date, time.Local)
	return err == nil && time.Now().Local().After(due)
}

var (
	// abbreviationMap maps shortcut keywords to date functions.
	abbreviationMap = map[string]func(time.Time) string{
		"td": formatToday, "tdy": formatToday, "today": formatToday,
		"tm": inDays(1), "tmmrw": inDays(1), "next": inDays(1),
		"af": inDays(2), "aft": inDays(2),
		"yd": inDays(-1), "yst": inDays(-1),
		"now": formatToday, "soon": inDays(3), "later": inDays(7),
		"someday": func(t time.Time) string { return "" },
		"nw":      inDays(7), "nxtwk": inDays(7),
		"n2w": inDays(14), "n3w": inDays(21),
		"eowk": nextWeekday(time.Friday),
		"nm":   inMonths(1), "em": endOfMonth,
		"mon": nextWeekday(time.Monday), "tue": nextWeekday(time.Tuesday),
		"wed": nextWeekday(time.Wednesday), "thu": nextWeekday(time.Thursday),
		"fri": nextWeekday(time.Friday), "sat": nextWeekday(time.Saturday),
		"sun":    nextWeekday(time.Sunday),
		"nxtmon": nextWeekday(time.Monday), "nxfri": nextWeekday(time.Friday),
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
	// monthNames maps month names/abbreviations to numbers.
	monthNames = map[string]int{
		"jan": 1, "january": 1, "feb": 2, "february": 2, "mar": 3, "march": 3,
		"apr": 4, "april": 4, "may": 5, "jun": 6, "june": 6, "jul": 7, "july": 7,
		"aug": 8, "august": 8, "sep": 9, "sept": 9, "september": 9,
		"oct": 10, "october": 10, "nov": 11, "november": 11, "dec": 12, "december": 12,
	}
	// relativeDateRegex matches formats like "2 days", "in 2 weeks", "next month".
	relativeDateRegex = regexp.MustCompile(`^(?:(?:in\s+)?|next\s+)?(\d+)?\s*(d|day|days|w|week|weeks|m|month|months)$`)
)

// ParseNaturalDate parses natural language dates into YYYY-MM-DD format.
// Examples: "today", "tomorrow", "fri", "in 2 days", "Dec 31, 2025", "2025-12-31".
func ParseNaturalDate(input string) (string, error) {
	input = strings.ToLower(strings.TrimSpace(input))
	if input == "" {
		return "", fmt.Errorf("date cannot be empty")
	}

	now := time.Now().Local()

	// Handle abbreviation shortcuts
	if f, ok := abbreviationMap[input]; ok {
		return f(now), nil
	}

	// Handle relative dates
	if matches := relativeDateRegex.FindStringSubmatch(input); matches != nil {
		unit := matches[2]
		num := 1
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

	// Handle natural language dates
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

// parseNaturalLanguageDate handles formats like "Dec 31, 2025" or "31 Dec 2025".
func parseNaturalLanguageDate(input string, now time.Time) (string, error) {
	parts := strings.Fields(input)
	if len(parts) < 2 || len(parts) > 3 {
		return "", fmt.Errorf("invalid natural language date")
	}

	var day, month, year int
	var err error

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

	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	if t.Year() != year || t.Month() != time.Month(month) || t.Day() != day {
		return "", fmt.Errorf("invalid date: %s", input)
	}
	return t.Format("2006-01-02"), nil
}

// ParseDateTimeDuration parses strings like "friday @ 18:00 for 1h" into date, time, and duration.
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

	date, err = ParseNaturalDate(main)
	if err != nil {
		return "", "", "", err
	}

	if timeStr != "" {
		if _, err := time.Parse("15:04", timeStr); err != nil {
			return "", "", "", fmt.Errorf("invalid time format: %s", timeStr)
		}
	}

	if duration != "" {
		if _, err := time.ParseDuration(duration); err != nil {
			return "", "", "", fmt.Errorf("invalid duration: %s", duration)
		}
	}

	return date, timeStr, duration, nil
}

// SetDueDate assigns a due date to a task by ID or text.
func SetDueDate(input string, dueDate string) error {
	tasks, err := LoadTasks()
	if err != nil {
		return fmt.Errorf("failed to load tasks: %v", err)
	}
	found := false
	parsedDate, err := ParseNaturalDate(dueDate)
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(input)
	for i, task := range tasks {
		if err == nil && task.ID == id {
			tasks[i].DueDate = parsedDate
			found = true
			break
		}
		if task.Text == input {
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

// formatToday returns the current date in YYYY-MM-DD format.
func formatToday(t time.Time) string {
	return t.Format("2006-01-02")
}

// inDays returns a function that adds n days to a date.
func inDays(n int) func(time.Time) string {
	return func(t time.Time) string {
		return t.AddDate(0, 0, n).Format("2006-01-02")
	}
}

// inMonths returns a function that adds n months to a date.
func inMonths(n int) func(time.Time) string {
	return func(t time.Time) string {
		return t.AddDate(0, n, 0).Format("2006-01-02")
	}
}

// endOfMonth returns the last day of the current month.
func endOfMonth(t time.Time) string {
	firstOfNext := time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, t.Location())
	return firstOfNext.AddDate(0, 0, -1).Format("2006-01-02")
}

// nextWeekday returns the next occurrence of the specified weekday.
func nextWeekday(wd time.Weekday) func(time.Time) string { // CalculateNextDueDate computes the next due date based on the recurrence rule.

	return func(t time.Time) string {
		offset := (int(wd) - int(t.Weekday()) + 7) % 7
		if offset == 0 {
			offset = 7
		}
		return t.AddDate(0, 0, offset).Format("2006-01-02")
	}
}

// startOfMonth returns the first day of the current month.
func startOfMonth(t time.Time) string {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).Format("2006-01-02")
}

// endOfWeek returns the last day of the current week (Saturday).
func endOfWeek(t time.Time) string {
	offset := 6 - int(t.Weekday())
	return t.AddDate(0, 0, offset).Format("2006-01-02")
}

// startOfNextMonth returns the first day of the next month.
func startOfNextMonth(t time.Time) string {
	next := t.AddDate(0, 1, 0)
	return time.Date(next.Year(), next.Month(), 1, 0, 0, 0, 0, t.Location()).Format("2006-01-02")
}

// endOfNextMonth returns the last day of the next month.
func endOfNextMonth(t time.Time) string {
	next := t.AddDate(0, 2, 0)
	return time.Date(next.Year(), next.Month(), 0, 0, 0, 0, 0, t.Location()).Format("2006-01-02")
}

// CalculateNextDueDate computes the next due date based on the recurrence rule.
func CalculateNextDueDate(currentDueDate, recurring string) string {
	if recurring == "" {
		return "" // No recurrence
	}

	now := time.Now().Local()
	currentDue, err := time.ParseInLocation("2006-01-02", currentDueDate, time.Local)
	if err != nil {
		currentDue = now // Fallback if invalid
	}

	recurring = strings.ToLower(strings.TrimSpace(recurring))
	if !strings.HasPrefix(recurring, "every ") {
		return "" // Invalid format
	}
	rule := strings.TrimPrefix(recurring, "every ")

	switch {
	case rule == "day" || rule == "daily":
		return currentDue.AddDate(0, 0, 1).Format("2006-01-02")
	case rule == "week" || rule == "weekly":
		return currentDue.AddDate(0, 0, 7).Format("2006-01-02")
	case rule == "month" || rule == "monthly":
		return currentDue.AddDate(0, 1, 0).Format("2006-01-02")
	case strings.Contains(rule, ","): // e.g., "mon,wed,fri"
		weekdays := strings.Split(rule, ",")
		var next time.Time = currentDue.AddDate(0, 0, 1) // Start from next day
		for {
			for _, wdStr := range weekdays {
				wd := parseWeekday(wdStr)
				if wd != -1 && next.Weekday() == time.Weekday(wd) {
					return next.Format("2006-01-02")
				}
			}
			next = next.AddDate(0, 0, 1)
		}
	default: // Single weekday, e.g., "sunday"
		wd := parseWeekday(rule)
		if wd == -1 {
			return "" // Invalid
		}
		return nextWeekday(time.Weekday(wd))(currentDue.AddDate(0, 0, 1)) // Next occurrence after current
	}
}

// parseWeekday maps day names to time.Weekday.
func parseWeekday(day string) int {
	dayMap := map[string]int{
		"mon": 1, "monday": 1,
		"tue": 2, "tuesday": 2,
		"wed": 3, "wednesday": 3,
		"thu": 4, "thursday": 4,
		"fri": 5, "friday": 5,
		"sat": 6, "saturday": 6,
		"sun": 0, "sunday": 0,
	}
	if val, ok := dayMap[day]; ok {
		return val
	}
	return -1
}
