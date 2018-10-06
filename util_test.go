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
	{"11,996", 11996},
	{"50%", 0.5},
	{"50 %", 0.5},
	{"01:00:00", 1},
	{"00:30:00", 0.5},
	{"30:36", 0.51},
	{"02:22", 0.0394444465637207},
	{"1 hour", 1},
	{"5 hours", 5},
	{"14 hours", 14},
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

func TestUtil_SanitizeKey(t *testing.T) {
	for _, tt := range keyTests {
		actual := SanitizeKey(tt.n)
		if actual != tt.expected {
			t.Errorf("SanitizeKey(%s): expected %s, actual %s", tt.n, tt.expected, actual)
		}
	}
}

func TestUtil_SanitizeKey_Parallel(t *testing.T) {
	for _, tt := range keyTests {
		go func(input string, output string) {
			actual := SanitizeKey(input)
			if actual != output {
				t.Errorf("SanitizeKey(%s): expected %s, actual %s", input, output, actual)
			}
		}(tt.n, tt.expected)
	}
}

func TestUtil_SanitizeValue(t *testing.T) {
	for _, tt := range valueTests {
		actual := SanitizeValue(tt.n)
		if actual != tt.expected {
			t.Errorf("SanitizeValue(%s): expected %f, actual %f", tt.n, tt.expected, actual)
		}
	}
}

func TestUtil_SanitizeValue_Parallel(t *testing.T) {
	for _, tt := range valueTests {
		go func(input string, output float32) {
			actual := SanitizeValue(input)
			if actual != output {
				t.Errorf("SanitizeValue(%s): expected %f, actual %f", input, output, actual)
			}
		}(tt.n, tt.expected)
	}
}

func TestUtil_Pluralizer(t *testing.T) {
	for _, tt := range pluralizerTests {
		actual := Pluralizer(tt.n)
		if actual != tt.expected {
			t.Errorf("Pluralizer(%s): expected %s, actual %s", tt.n, tt.expected, actual)
		}
	}
}

func TestUtil_Pluralizer_Parallel(t *testing.T) {
	for _, tt := range pluralizerTests {
		go func(input string, output string) {
			actual := Pluralizer(input)
			if actual != output {
				t.Errorf("Pluralizer(%s): expected %s, actual %s", input, output, actual)
			}
		}(tt.n, tt.expected)
	}
}
