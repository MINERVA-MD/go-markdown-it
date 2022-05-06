package pkg

import (
	"strings"
)

func Normalize(state *StateCore, _ *StateBlock, _ *StateInline, _ int, _ int, _ bool) bool {
	var src string

	// Normalize newlines
	src = NEWLINES_RE.ReplaceAllString(state.Src, "\n")

	// Replace NULL characters
	src = strings.Replace(src, "\x00", "\uFFFD", -1)

	state.Src = src

	return true
}
