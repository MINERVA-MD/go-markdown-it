package core

import (
	. "go-markdown-it/pkg/common"
	. "go-markdown-it/pkg/types"
	"strings"
)

func Normalize(state *StateCore) {
	var src string

	// Normalize newlines
	src = NEWLINES_RE.ReplaceAllString(src, `\n`)

	// Replace NULL characters
	src = strings.Replace(src, "\x00", "\uFFFD", -1)

	state.Src = src
}
