package tools

import (
	"strings"
	"unicode/utf8"
)

type PadDirection int

const (
	LEFT PadDirection = iota
	CENTER
	RIGHT
)

// Pad the string content
func StrPad(direction PadDirection, text string, width int) string {
	output := ""
	textwith := utf8.RuneCountInString(text)
	nbspaces := width - textwith

	if nbspaces <= 0 {
		return text
	}

	switch direction {
	case RIGHT:
		output = strings.Repeat(" ", nbspaces) + text
	case LEFT:
		output = text + strings.Repeat(" ", nbspaces)
	case CENTER:
		leftspaces := int(nbspaces / 2)
		rightspace := nbspaces - leftspaces

		output = strings.Repeat(" ", leftspaces) + text + strings.Repeat(" ", rightspace)
	}

	return output
}
