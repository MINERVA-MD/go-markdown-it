package core

import (
	. "go-markdown-it/pkg/common"
	. "go-markdown-it/pkg/rules/block"
	. "go-markdown-it/pkg/rules/inline"
	"strings"
)

func Normalize(state *StateCore, _ *StateBlock, _ *StateInline, _ int, _ int, _ bool) bool {
	var src string

	// Normalize newlines
	src = NEWLINES_RE.ReplaceAllString(src, `\n`)

	// Replace NULL characters
	src = strings.Replace(src, "\x00", "\uFFFD", -1)

	state.Src = src

	return true
}
