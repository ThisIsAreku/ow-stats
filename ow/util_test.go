package ow

import "testing"

var valueTests = []struct {
	n        string // input
	expected string // expected result
}{
	{"Multikill - Best", "multikill_best"},
	{"Soldier: 76", "soldier76"},
	{"Lúcio", "lucio"},
}

func TestSanitizeKey(t *testing.T) {
	for _, tt := range valueTests {
		actual := SanitizeKey(tt.n)
		if actual != tt.expected {
			t.Errorf("SanitizeKey(%s): expected %s, actual %s", tt.n, tt.expected, actual)
		}
	}
}
