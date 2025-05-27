package todo

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// 📌 Centralized abbreviation map — natural language → date string
var abbreviationMap = map[string]func(time.Time) string{
	// 🗓️ Absolute
	"td": formatToday, "tdy": formatToday, "today": formatToday,
	"tm": inDays(1), "tmmrw": inDays(1), "next": inDays(1),
	"af": inDays(2), "aft": inDays(2),
	"yd": inDays(-1), "yst": inDays(-1),
	"now": formatToday, "soon": inDays(3), "later": inDays(7),
	"someday": func(t time.Time) string { return "" },

	// 🗓️ Weeks
	"nw": inDays(7), "nxtwk": inDays(7),
	"n2w": inDays(14), "n3w": inDays(21),
	"eowk": nextWeekday(time.Friday),

	// 🗓️ Months
	"nm": inMonths(1), "em": endOfMonth,

	// 🗓️ Weekdays (next)
	"mon": nextWeekday(time.Monday), "tue": nextWeekday(time.Tuesday),
	"wed": nextWeekday(time.Wednesday), "thu": nextWeekday(time.Thursday),
	"fri": nextWeekday(time.Friday), "sat": nextWeekday(time.Saturday),
	"sun": nextWeekday(time.Sunday),

	"nxtmon": nextWeekday(time.Monday), "nxfri": nextWeekday(time.Friday),

	// 🕓 Misc
	"eod": formatToday,
	"ew": func(t time.Time) string {
		return t.AddDate(0, 0, 7-int(t.Weekday())).Format("2006-01-02")
	},
}

// 📆 Parses strings like "in 3 days", "fri", "2024-04-25"
func ParseNaturalDate(input string) (string, error) {
	input = strings.ToLower(strings.TrimSpace(input))
	now := time.Now()

	// Shorthand: "tomorrow", "fri", etc.
	if f, ok := abbreviationMap[input]; ok {
		return f(now), nil
	}

	// Handle "in N units"
	if strings.HasPrefix(input, "in ") {
		parts := strings.Fields(input[3:])
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid relative date format: %s", input)
		}

		num, err := strconv.Atoi(parts[0])
		if err != nil {
			return "", fmt.Errorf("invalid number in relative date: %v", err)
		}

		unit := parts[1]
		switch unit {
		case "d", "day", "days":
			return now.AddDate(0, 0, num).Format("2006-01-02"), nil
		case "w", "week", "weeks":
			return now.AddDate(0, 0, num*7).Format("2006-01-02"), nil
		case "m", "month", "months":
			return now.AddDate(0, num, 0).Format("2006-01-02"), nil
		default:
			return "", fmt.Errorf("unsupported unit: %s", unit)
		}
	}

	// Fallback: date strings like "2025-05-02"
	for _, layout := range []string{"2006-01-02", "02-01-2006"} {
		if t, err := time.Parse(layout, input); err == nil {
			return t.Format("2006-01-02"), nil
		}
	}

	return "", fmt.Errorf("could not parse date: %s", input)
}

// ParseDateTimeDuration parses a string like "friday @ 18:00 for 1h"
func ParseDateTimeDuration(input string) (date string, timeStr string, duration string, err error) {
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

	// validate HH:MM
	if timeStr != "" {
		if _, err := time.Parse("15:04", timeStr); err != nil {
			return "", "", "", fmt.Errorf("invalid time format: %s", timeStr)
		}
	}

	// validate duration (like 30m, 1h, 90m)
	if duration != "" {
		if _, err := time.ParseDuration(duration); err != nil {
			return "", "", "", fmt.Errorf("invalid duration: %s", duration)
		}
	}

	return date, timeStr, duration, nil
}

// 📆 Helpers
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

func formatToday(t time.Time) string {
	return t.Format("2006-01-02")
}
