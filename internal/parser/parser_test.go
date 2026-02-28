package parser

import (
	"testing"
)

func TestParseSet(t *testing.T) {
	tests := []struct {
		input    string
		wantErr  bool
		wantSecs float64
	}{
		// Valid inputs
		{"set 90", false, 90},
		{"set 1:30", false, 90},
		{"set 0", false, 0},
		{"set 2:05", false, 125},
		{"set 0:00", false, 0},
		{"set 10:00", false, 600},

		// Missing keyword
		{"90", true, 0},
		{"1:30", true, 0},

		// Missing duration
		{"set", true, 0},

		// Non-numeric
		{"set abc", true, 0},
		{"set 1:xx", true, 0},

		// Seconds out of range
		{"set 1:60", true, 0},
		{"set 1:99", true, 0},

		// Negative values
		{"set -1", true, 0},
		{"set -1:30", true, 0},

		// Too many colons
		{"set 1:2:3", true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p, err := ParseSet(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			got := p.TimeDisplay().Seconds()
			if got != tt.wantSecs {
				t.Errorf("got %.0f seconds, want %.0f", got, tt.wantSecs)
			}
		})
	}
}
