package block

import (
	. "go-markdown-it/pkg/common"
	"strings"
)

func (state *StateBlock) Hr(startLine int, endLine int, silent bool) bool {

	pos := state.BMarks[startLine] + state.TShift[startLine]
	max := state.EMarks[startLine]

	if state.SCount[startLine]-state.BlkIndent >= 4 {
		return false
	}

	marker := CharCodeAt(state.Src, pos)
	pos++

	// Check hr marker
	if marker != 0x2A /* * */ &&
		marker != 0x2D /* - */ &&
		marker != 0x5F /* _ */ {
		return false
	}

	// markers can be mixed with spaces, but there should be at least 3 of them

	cnt := 1

	for pos < max {
		ch := CharCodeAt(state.Src, pos)
		pos++

		if ch != marker && !IsSpace(ch) {
			return false
		}
		if ch == marker {
			cnt++
		}
	}

	if cnt < 3 {
		return false
	}

	if silent {
		return true
	}

	state.Line = startLine + 1

	token := state.Push("hr", "hr", 0)
	token.Map = []int{startLine, state.Line}
	token.Markup = strings.Join(make([]string, cnt+1)[:], string(marker))

	return true
}
