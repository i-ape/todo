package config

import (
	"os"
	"testing"
)

func TestGetDefaultPriority(t *testing.T) {
	defer os.Unsetenv("TODO_DEFAULT_PRIORITY")

	tests := []struct {
		name string
		env  string
		want string
	}{
		{"default", "", "medium"},
		{"high", "high", "high"},
		{"low", "low", "low"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env != "" {
				os.Setenv("TODO_DEFAULT_PRIORITY", tt.env)
			} else {
				os.Unsetenv("TODO_DEFAULT_PRIORITY")
			}
			got := GetDefaultPriority()
			if got != tt.want {
				t.Errorf("GetDefaultPriority() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDisableFzf(t *testing.T) {
	defer os.Unsetenv("TODO_DISABLE_FZF")

	tests := []struct {
		name string
		env  string
		want bool
	}{
		{"default", "", false},
		{"true", "true", true},
		{"false", "false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env != "" {
				os.Setenv("TODO_DISABLE_FZF", tt.env)
			} else {
				os.Unsetenv("TODO_DISABLE_FZF")
			}
			got := DisableFzf
			if got != tt.want {
				t.Errorf("DisableFzf = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSortOrder(t *testing.T) {
	defer os.Unsetenv("TODO_SORT_ORDER")

	tests := []struct {
		name string
		env  string
		want string
	}{
		{"default", "", "due"},
		{"priority", "priority", "priority"},
		{"text", "text", "text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.env != "" {
				os.Setenv("TODO_SORT_ORDER", tt.env)
			} else {
				os.Unsetenv("TODO_SORT_ORDER")
			}
			got := GetSortOrder()
			if got != tt.want {
				t.Errorf("GetSortOrder() = %q, want %q", got, tt.want)
			}
		})
	}
}
