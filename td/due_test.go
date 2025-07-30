package todo

import (
	"testing"
	"time"
)

func TestParseNaturalDate(t *testing.T) {
	now := time.Now().Local()
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"today", now.Format("2006-01-02"), false},
		{"tomorrow", now.AddDate(0, 0, 1).Format("2006-01-02"), false},
		{"fri", nextWeekday(time.Friday)(now), false},
		{"in 2 days", now.AddDate(0, 0, 2).Format("2006-01-02"), false},
		{"next week", now.AddDate(0, 0, 7).Format("2006-01-02"), false},
		{"Dec 31, 2025", "2025-12-31", false},
		{"31 Dec 2025", "2025-12-31", false},
		{"2025-12-31", "2025-12-31", false},
		{"31-12-2025", "2025-12-31", false},
		{"12-31-2025", "2025-12-31", false},
		{"invalid", "", true},
		{"", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseNaturalDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseNaturalDate(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("ParseNaturalDate(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIsOverdue(t *testing.T) {
	now := time.Now().Local()
	tests := []struct {
		date     string
		expected bool
	}{
		{now.AddDate(0, 0, -1).Format("2006-01-02"), true}, // Yesterday
		{now.Format("2006-01-02"), false},                  // Today
		{now.AddDate(0, 0, 1).Format("2006-01-02"), false}, // Tomorrow
		{"invalid", false},                                 // Invalid date
	}
	for _, tt := range tests {
		t.Run(tt.date, func(t *testing.T) {
			got := IsOverdue(tt.date)
			if got != tt.expected {
				t.Errorf("IsOverdue(%q) = %v, want %v", tt.date, got, tt.expected)
			}
		})
	}
}

func TestParseDateTimeDuration(t *testing.T) {
	tests := []struct {
		input    string
		wantDate string
		wantTime string
		wantDur  string
		wantErr  bool
	}{
		{"friday @ 18:00 for 1h", nextWeekday(time.Friday)(time.Now().Local()), "18:00", "1h", false},
		{"tomorrow @ 10:30", time.Now().Local().AddDate(0, 0, 1).Format("2006-01-02"), "10:30", "", false},
		{"today", time.Now().Local().Format("2006-01-02"), "", "", false},
		{"invalid @ 25:00", "", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			date, timeStr, dur, err := ParseDateTimeDuration(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDateTimeDuration(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if date != tt.wantDate {
				t.Errorf("ParseDateTimeDuration(%q) date = %q, want %q", tt.input, date, tt.wantDate)
			}
			if timeStr != tt.wantTime {
				t.Errorf("ParseDateTimeDuration(%q) time = %q, want %q", tt.input, timeStr, tt.wantTime)
			}
			if dur != tt.wantDur {
				t.Errorf("ParseDateTimeDuration(%q) duration = %q, want %q", tt.input, dur, tt.wantDur)
			}
		})
	}
}

func TestParseNaturalLanguageDate(t *testing.T) {
	now := time.Now().Local()
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"Dec 31, 2025", "2025-12-31", false},
		{"31 Dec 2025", "2025-12-31", false},
		{"31 Dec", now.Format("2006") + "-12-31", false},
		{"invalid month", "", true},
		{"31 invalid 2025", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseNaturalLanguageDate(tt.input, now)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseNaturalLanguageDate(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.expected {
				t.Errorf("parseNaturalLanguageDate(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
