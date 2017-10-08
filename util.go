package ow_stats

import (
	"regexp"
	"strings"
	"strconv"
	"log"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/runes"
	"unicode"
)

var spaceRegex = regexp.MustCompile(`[-\s]`)
var nonAlphaRegex = regexp.MustCompile(`\W`)
var normalizeRegex = regexp.MustCompile(`_{2,}`)
var hourRegex = regexp.MustCompile(`(?P<Val>[0-9]+) hours?`)
var minuteRegex = regexp.MustCompile(`(?P<Val>[0-9]+) minutes?`)
var secondRegex = regexp.MustCompile(`(?P<Val>[0-9]+(?:\.[0-9]+)?) seconds?`)
var percentRegex = regexp.MustCompile(`(?P<Val>[0-9]{1,3})\s?%`)
var pluralizerRegex = regexp.MustCompile(`(?P<Start>_|[^a-z]|^)(?P<Term>blow|boost|kill|assist|barrier|hit|multikill|elimination)(?P<End>$|[^a-z])`)

var t = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

func SanitizeKey(text string) string {
	text = strings.ToLower(text)
	text, _, _ = transform.String(t, text)
	text = spaceRegex.ReplaceAllString(text, "_")
	text = nonAlphaRegex.ReplaceAllString(text, "")
	text = normalizeRegex.ReplaceAllString(text, "_")

	text = strings.Replace(text, "soldier_76", "soldier76", 1)

	return text
}

func SanitizeValue(text string) float32 {
	text = strings.ToLower(text)

	fnParseFloat := func(text string) (float32, error) {
		longNum := strings.Replace(text, ",", "", -1)
		v, err := strconv.ParseFloat(longNum, 32)
		if err != nil {
			return 0., err
		}

		return float32(v), nil
	}

	fnParseSimpleFloat32 := func(t string) float32 {
		if v, err := strconv.Atoi(t); err == nil {
			return float32(v)
		}

		return 0
	}

	if text == "--" {
		return .0
	}

	if v, err := fnParseFloat(text); err == nil {
		return float32(v)
	}

	// value is a percentage
	if m := percentRegex.FindStringSubmatch(text); m != nil {
		if v, err := fnParseFloat(m[1]); err == nil {
			return v / 100.
		}

		log.Printf("Value is a percentage but failed to decode: %s\n", text)

		return 0
	}

	if strings.ContainsRune(text, ':') {
		parts := strings.Split(text, ":")

		var h, m, s float32
		switch len(parts) {
		case 3:
			h, m, s = fnParseSimpleFloat32(parts[0]), fnParseSimpleFloat32(parts[1]), fnParseSimpleFloat32(parts[2])
		case 2:
			m, s = fnParseSimpleFloat32(parts[0]), fnParseSimpleFloat32(parts[1])
		case 1:
			s = fnParseSimpleFloat32(parts[0])
		}

		return h + ((m + (s / 60.0)) / 60.0)
	}

	if m := hourRegex.FindStringSubmatch(text); m != nil {
		if v, err := fnParseFloat(m[1]); err == nil {
			return v
		}

		log.Printf("Value is an hour but failed to decode: %s\n", text)

		return 0
	}

	if m := minuteRegex.FindStringSubmatch(text); m != nil {
		if v, err := fnParseFloat(m[1]); err == nil {
			return v / 60.0
		}

		log.Printf("Value is a minute but failed to decode: %s\n", text)

		return 0
	}

	if m := secondRegex.FindStringSubmatch(text); m != nil {
		if v, err := fnParseFloat(m[1]); err == nil {
			return v / 3600.0
		}

		log.Printf("Value is a second but failed to decode: %s\n", text)

		return 0

	}

	log.Printf("Unable to find anything to parse: %s\n", text)

	return 0
}

func Pluralizer(text string) string {
	text = strings.ToLower(text)
	text = pluralizerRegex.ReplaceAllString(text, `${1}${2}s${3}`)

	text = strings.Replace(text, "kills_streak_", "kill_streak_", 1)

	return text
}
