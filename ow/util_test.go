package ow

import "testing"

var keyTests = []struct {
	n        string // input
	expected string // expected result
}{
	{"Multikill - Best", "multikill_best"},
	{"Soldier: 76", "soldier76"},
	{"Lúcio", "lucio"},
}

var valueTests = []struct {
	n        string  // input
	expected float32 // expected result
}{
	{"5", 5},
	{"50%", 0.5},
	{"50 %", 0.5},
	{"01:00:00", 1},
	{"00:30:00", 0.5},
	{"30:36", 0.51},
	{"02:22", 0.0394444465637207},
	{"1 hour", 1},
	{"5 hours", 5},
	{"30 minutes", 0.5},
	{"2.8", 2.8},
}

func TestSanitizeKey(t *testing.T) {
	for _, tt := range keyTests {
		actual := SanitizeKey(tt.n)
		if actual != tt.expected {
			t.Errorf("SanitizeKey(%s): expected %s, actual %s", tt.n, tt.expected, actual)
		}
	}
}

func TestSanitizeValue(t *testing.T) {
	for _, tt := range valueTests {
		actual := SanitizeValue(tt.n)
		if actual != tt.expected {
			t.Errorf("SanitizeValue(%s): expected %f, actual %f", tt.n, tt.expected, actual)
		}
	}
}
