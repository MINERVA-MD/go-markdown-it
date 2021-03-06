package pkg

import (
	"strings"
)

func Fence(
	_ *StateCore,
	state *StateBlock,
	_ *StateInline,
	startLine int,
	endLine int,
	silent bool,
) bool {
	return state.Fence(startLine, endLine, silent)
}

func (state *StateBlock) Fence(startLine int, endLine int, silent bool) bool {

	var mem int
	var marker rune
	var markup string
	var params string
	haveEndMarker := false
	pos := state.BMarks[startLine] + state.TShift[startLine]
	max := state.EMarks[startLine]

	// if it's indented more than 3 spaces, it should be a code block
	if state.SCount[startLine]-state.BlkIndent >= 4 {
		return false
	}

	if pos+3 > max {
		return false
	}

	marker, _ = state.Src2.CharCodeAt(pos)

	if marker != 0x7E /* ~ */ && marker != 0x60 /* ` */ {
		//fmt.Println("Returning false")
		return false
	}

	// scan marker length
	mem = pos
	pos = state.SkipChars(pos, marker)

	_len := pos - mem

	if _len < 3 {
		return false
	}

	markup = state.Src2.Slice(mem, pos)
	params = state.Src2.Slice(pos, max)

	if marker == 0x60 {
		if strings.Contains(params, string(marker)) {
			return false
		}
	}

	// Since start is found, we can report success here in validation mode
	if silent {
		return true
	}

	// search end of block
	nextLine := startLine

	for {
		nextLine++
		if nextLine >= endLine {
			// unclosed block should be autoclosed by end of document.
			// also block seems to be autoclosed by end of parent
			break
		}

		mem = state.BMarks[nextLine] + state.TShift[nextLine]
		pos = state.BMarks[nextLine] + state.TShift[nextLine]
		max = state.EMarks[nextLine]

		if pos < max && state.SCount[nextLine] < state.BlkIndent {
			// non-empty line with negative indent should stop the list:
			// - ```
			//  test
			break
		}

		if cc, _ := state.Src2.CharCodeAt(pos); cc != marker {
			continue
		}

		if state.SCount[nextLine]-state.BlkIndent >= 4 {
			// closing fence should be indented less than 4 spaces
			continue
		}

		pos = state.SkipChars(pos, marker)

		// closing code fence must be at least as long as the opening one
		if (pos - mem) < _len {
			continue
		}

		// make sure tail has spaces only
		pos = state.SkipSpaces(pos)

		if pos < max {
			continue
		}

		haveEndMarker = true
		// found!
		break
	}

	// If a fence has heading spaces, they should be removed from its inner block
	_len = state.SCount[startLine]

	state.Line = nextLine

	if haveEndMarker {
		state.Line += 1
	}

	token := state.Push("fence", "code", 0)
	token.Info = params
	token.Content = state.GetLines(startLine+1, nextLine, _len, true)
	token.Markup = markup
	token.Map = []int{startLine, state.Line}

	return true
}
