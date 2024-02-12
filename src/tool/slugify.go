package tool

import (
	"regexp"
	"strings"
	"unicode"
)

func Slugify(s string) string {
	var slug string
	isAllowed := func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			return r
		}
		if unicode.IsSpace(r) {
			return '-'
		}
		return -1
	}

	slug = strings.Map(isAllowed, strings.ToLower(s))

	// Remove any double dashes caused by spaces next to disallowed chars
	reg, _ := regexp.Compile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	return slug
}
