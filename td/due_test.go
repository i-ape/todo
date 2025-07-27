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
