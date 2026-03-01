package parser

import (
	"testing"
	"time"

	"github.com/BobbyGerace/workout-timer/internal/types"
)

func TestParseSet(t *testing.T) {
	auto := types.ModeAuto
	manual := types.ModeManual

	tests := []struct {
		input       string
		defaultMode types.Mode
		wantErr     bool
		// expected first interval duration in seconds
		wantFirstSecs float64
		// expected number of intervals
		wantIntervals int
		// expected rounds (0 = loop forever)
		wantRounds int
	}{
		// ── Single duration, no flags ──────────────────────────────────────
		{"set 90", auto, false, 90, 1, 0},
		{"set 1:30", auto, false, 90, 1, 0},
		{"set 0", auto, false, 0, 1, 0},
		{"set 10:00", auto, false, 600, 1, 0},

		// ── Mode flag ─────────────────────────────────────────────────────
		{"set auto 60", auto, false, 60, 1, 0},
		{"set manual 60", auto, false, 60, 1, 0},

		// ── Default mode falls through ─────────────────────────────────────
		{"set 60", manual, false, 60, 1, 0},

		// ── Round counts ──────────────────────────────────────────────────
		{"set 60 x5", auto, false, 60, 1, 5},
		{"set auto 60 x10", auto, false, 60, 1, 10},
		{"set manual 60 x1", auto, false, 60, 1, 1},

		// ── Multi-interval ────────────────────────────────────────────────
		{"set 90,60", auto, false, 90, 2, 0},
		{"set auto 1:30,60,4:00", auto, false, 90, 3, 0},
		{"set auto 1:30,60,4:00 x3", auto, false, 90, 3, 3},

		// ── Errors: missing keyword / duration ────────────────────────────
		{"90", auto, true, 0, 0, 0},
		{"set", auto, true, 0, 0, 0},
		{"set auto", auto, true, 0, 0, 0},

		// ── Errors: invalid duration ──────────────────────────────────────
		{"set abc", auto, true, 0, 0, 0},
		{"set 1:60", auto, true, 0, 0, 0},
		{"set -1", auto, true, 0, 0, 0},
		{"set 1:2:3", auto, true, 0, 0, 0},
		{"set 90,abc", auto, true, 0, 0, 0},

		// ── Errors: bad round count ────────────────────────────────────────
		{"set 60 x0", auto, true, 0, 0, 0},
		{"set 60 x-1", auto, true, 0, 0, 0},
		{"set 60 xfoo", auto, true, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p, err := ParseSet(tt.input, tt.defaultMode)
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
			if got != tt.wantFirstSecs {
				t.Errorf("first interval: got %.0fs, want %.0fs", got, tt.wantFirstSecs)
			}

			// Verify interval count via IntervalProgress (or check total == 1 for single)
			cur, total := p.IntervalProgress()
			if tt.wantIntervals > 1 {
				if total != tt.wantIntervals {
					t.Errorf("interval total: got %d, want %d", total, tt.wantIntervals)
				}
				if cur != 1 {
					t.Errorf("interval current: got %d, want 1", cur)
				}
			} else {
				// Single interval — IntervalProgress returns (0,0)
				if cur != 0 || total != 0 {
					t.Errorf("single interval: expected IntervalProgress (0,0), got (%d,%d)", cur, total)
				}
			}

			// Verify round count via RoundProgress
			_, roundTotal := p.RoundProgress()
			if tt.wantRounds > 1 {
				if roundTotal != tt.wantRounds {
					t.Errorf("round total: got %d, want %d", roundTotal, tt.wantRounds)
				}
			} else {
				if roundTotal != 0 {
					t.Errorf("expected RoundProgress total 0, got %d", roundTotal)
				}
			}
		})
	}
}

func TestParseDurationList(t *testing.T) {
	tests := []struct {
		input    string
		wantErr  bool
		wantDurs []time.Duration
	}{
		{"90", false, []time.Duration{90 * time.Second}},
		{"1:30,60", false, []time.Duration{90 * time.Second, 60 * time.Second}},
		{"1:30,60,4:00", false, []time.Duration{90 * time.Second, 60 * time.Second, 240 * time.Second}},
		{"abc", true, nil},
		{"90,bad", true, nil},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseDurationList(tt.input)
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
			if len(got) != len(tt.wantDurs) {
				t.Errorf("got %d durations, want %d", len(got), len(tt.wantDurs))
				return
			}
			for i, d := range got {
				if d != tt.wantDurs[i] {
					t.Errorf("[%d] got %v, want %v", i, d, tt.wantDurs[i])
				}
			}
		})
	}
}

func TestParseRounds(t *testing.T) {
	tests := []struct {
		input   string
		wantErr bool
		want    int
	}{
		{"x1", false, 1},
		{"x5", false, 5},
		{"x10", false, 10},
		{"x0", true, 0},
		{"x-1", true, 0},
		{"xfoo", true, 0},
		{"5", true, 0},
		{"", true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseRounds(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for %q, got nil", tt.input)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}
