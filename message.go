package sos

import (
	"fmt"
	"strings"
)

// FallbackMessage creates a simple error message as the default error message.
func FallbackMessage(code Code) string {
	return fmt.Sprintf("%s error", code)
}

var (
	replacer = *strings.NewReplacer(
		"<!!P!!>", "%%",
		"%%", "<!!P!!>",
		"%d", "%v",
		"%f", "%v",
		"%s", "%v",
		"%t", "%v",
	)
)

func sprintf(s string, args ...interface{}) string {

	if len(args) == 0 {
		return s // If no args provided then just return the string param.
	}

	// Get the number of placeholder values which aren't literal '%'.
	s = replacer.Replace(s)
	count := strings.Count(s, "%")
	s = replacer.Replace(s)

	if len(args) == count {
		return fmt.Sprintf(s, args...)
	}

	if count < len(args) {
		// Since placeholder count is less than args then just take that number of args.
		return fmt.Sprintf(s, args[:count]...)
	}

	// Now try to provide cleaner output than fmt.Sprintf would.

	var b strings.Builder

	var found int
	num := count - len(args)

	for _, r := range s {
		if r == '%' {
			if found == num {
				break // All available placeholders have been consumed.
			}
			found++
		}
		b.WriteRune(r)
	}

	// Adding elipsis to indicate that some placeholder values were ignored.

	s = fmt.Sprintf("%s...", strings.TrimSpace(b.String()))

	return fmt.Sprintf(s, args...)
}
