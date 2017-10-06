package ow

import (
	"regexp"
	"strings"
	"strconv"
	"log"
)

var spaceRegex = regexp.MustCompile(`[-\s]`)
var nonAlphaRegex = regexp.MustCompile(`\W`)
var normalizeRegex = regexp.MustCompile(`_{2,}`)

func SanitizeKey(text string) string {
	text = strings.ToLower(text)
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
	if text[len(text)-1:] == "%" {
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

		return h + ((m + (s / 60)) / 60)
	}

	if strings.HasSuffix(text, " hours") {
		if v, err := strconv.Atoi(text[:len(text)-6]); err == nil {
			return float32(v)
		}

		log.Printf("Value is an hour but failed to decode: %s\n", text)

		return 0
	}

	if strings.HasSuffix(text, " minutes") {
		if v, err := strconv.Atoi(text[:len(text)-6]); err == nil {
			return float32(v) / 60.0
		}

		log.Printf("Value is a minute but failed to decode: %s\n", text)

		return 0
	}

	if strings.HasSuffix(text, " seconds") {
		if v, err := strconv.Atoi(text[:len(text)-6]); err == nil {
			return float32(v) / 3600.0
		}

		log.Printf("Value is a second but failed to decode: %s\n", text)

		return 0

	}

	return 0
}
