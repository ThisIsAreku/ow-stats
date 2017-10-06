package ow

import (
	"regexp"
	"strings"
	"strconv"
	"log"
	"fmt"
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
var secondRegex = regexp.MustCompile(`(?P<Val>[0-9]+\.?[0-9]+) seconds?`)
var percentRegex = regexp.MustCompile(`(?P<Val>[0-9]{1,3})\s?\%`)

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

	longNum := strings.Replace(text, ",", "", -1)
	if v, err := strconv.Atoi(longNum); err == nil {
		return float32(v)
	}

	// value is a percentage
	if m := percentRegex.FindAllString(text, -1); len(m) > 0 {
		fmt.Printf("TTTT: %+v\n", m)
		if v, err := strconv.Atoi(text[:len(text)-1]); err == nil {
			return float32(v) / 100
		}

		log.Printf("Value is a percentage but failed to decode: %s\n", text)

		return 0
	}

	if strings.ContainsRune(text, ':') {
		parts := strings.Split(text, ":")
		fnToSimpleFloat32 := func(t string) float32 {
			if v, err := strconv.Atoi(t); err == nil {
				return float32(v)
			}

			return 0
		}

		var h, m, s float32
		switch len(parts) {
		case 3:
			h, m, s = fnToSimpleFloat32(parts[0]), fnToSimpleFloat32(parts[1]), fnToSimpleFloat32(parts[2])
		case 2:
			m, s = fnToSimpleFloat32(parts[0]), fnToSimpleFloat32(parts[1])
		case 1:
			s = fnToSimpleFloat32(parts[0])
		}

		return h + ((m + (s / 60.0)) / 60.0)
	}

	if m := hourRegex.FindStringSubmatch(text); m != nil {
		if v, err := strconv.Atoi(m[1]); err == nil {
			return float32(v)
		}

		log.Printf("Value is an hour but failed to decode: %s\n", text)

		return 0
	}

	if m := minuteRegex.FindStringSubmatch(text); m != nil {
		if v, err := strconv.Atoi(m[1]); err == nil {
			return float32(v) / 60.0
		}

		log.Printf("Value is a minute but failed to decode: %s\n", text)

		return 0
	}

	if m := secondRegex.FindStringSubmatch(text); m != nil {
		if v, err := strconv.Atoi(m[1]); err == nil {
			return float32(v) / 3600.0
		}

		log.Printf("Value is a second but failed to decode: %s\n", text)

		return 0

	}

	return 0
}
