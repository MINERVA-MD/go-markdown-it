package block

import (
	. "go-markdown-it/pkg"
	. "go-markdown-it/pkg/common"
	"strconv"
	"strings"
)

func Heading(state *StateBlock, startLine int, endLine int, silent bool) bool {
	pos := state.BMarks[startLine] + state.TShift[startLine]
	max := state.EMarks[startLine]

	// if it's indented more than 3 spaces, it should be a code block
	if state.SCount[startLine]-state.BlkIndent >= 4 {
		return false
	}

	ch := CharCodeAt(state.Src, pos)

	if ch != 0x23 /* # */ || pos >= max {
		return false
	}

	level := 1
	pos++
	ch = CharCodeAt(state.Src, pos)

	for ch == 0x23 /* # */ && pos < max && level <= 6 {
		level++
		pos++
		ch = CharCodeAt(state.Src, pos)
	}

	if level > 6 || (pos < max && !IsSpace(ch)) {
		return false
	}

	if silent {
		return true
	}

	// Let's cut tails like '    ###  ' from the end of string

	max = state.SkipSpacesBack(max, pos)
	tmp := state.SkipCharsBack(max, 0x23, pos)

	if tmp > pos && IsSpace(CharCodeAt(state.Src, tmp-1)) {
		max = tmp
	}

	state.Line = startLine + 1

	// Opening Heading
	token := state.Push("heading_open", "h"+strconv.Itoa(level), 1)
	token.Markup = "######"[0:level]
	token.Map = []int{startLine, state.Line}

	//  Contents
	token = state.Push("inline", "", 0)
	token.Content = strings.TrimSpace(string([]rune(state.Src)[pos:max]))
	token.Map = []int{startLine, state.Line}
	token.Children = []*Token{}

	// Closing Heading
	token = state.Push("heading_close", "h"+strconv.Itoa(level), -1)
	token.Markup = "######"[0:level]

	return true
}
