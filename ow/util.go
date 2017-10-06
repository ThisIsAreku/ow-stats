package ow

import (
	"regexp"
	"strings"
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
