package ow_stats

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
	{"1 minute", 1 / 60.},
	{"50 seconds", 50 / 3600.},
	{"2.8 seconds", 2.8 / 3600.},
	{"1 second", 1 / 3600.},
	{"2.8", 2.8},
}

var pluralizerTests = []struct {
	n        string // input
	expected string // expected result
}{
	{"solo_kill", "solo_kills"},
	{"solo_kills", "solo_kills"},
	{"final_blow", "final_blows"},
	{"final_blows", "final_blows"},
	{"nano_boost_applied", "nano_boosts_applied"},
	{"nano_boosts_applied", "nano_boosts_applied"},
	{"final_blows_most_in_game", "final_blows_most_in_game"},
	{"final_blow_most_in_game", "final_blows_most_in_game"},
	{"projected_barriers_applied", "projected_barriers_applied"},
	{"projected_barrier_applied", "projected_barriers_applied"},
	{"multikills", "multikills"},
	{"multikill", "multikills"},
	{"eliminations_most_in_game", "eliminations_most_in_game"},
	{"elimination_most_in_game", "eliminations_most_in_game"},
	{"testblow_test", "testblow_test"},
	{"kill_streak_best", "kill_streak_best"},
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

func TestPluralizer(t *testing.T) {
	for _, tt := range pluralizerTests {
		actual := Pluralizer(tt.n)
		if actual != tt.expected {
			t.Errorf("Pluralizer(%s): expected %s, actual %s", tt.n, tt.expected, actual)
		}
	}
}
