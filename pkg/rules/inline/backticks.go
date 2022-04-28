package inline

import (
	. "go-markdown-it/pkg/common"
	. "go-markdown-it/pkg/rules/block"
	"strings"
)

func (state *StateInline) BackLinks(silent bool) bool {
	pos := state.Pos
	ch := CharCodeAt(state.Src, pos)

	if ch != 0x60 /* ` */ {
		return false
	}

	start := pos
	pos++
	max := state.PosMax

	// scan marker length
	for pos < max && CharCodeAt(state.Src, pos) == 0x60 /* ` */ {
		pos++
	}

	marker := state.Src[start:pos]
	openerLength := len(marker)

	backTicksCheck := 0

	if state.Backticks[openerLength] != 0 {
		backTicksCheck = state.Backticks[openerLength]
	}

	if state.BackTicksScanned && (backTicksCheck <= start) {
		if !silent {
			state.Pending += marker
		}

		state.Pos += openerLength
		return true
	}

	matchStart := pos
	matchEnd := pos

	// Nothing found in the cache, scan until the end of the line (or until marker is found)
	for {
		matchStart = strings.Index(state.Src[matchEnd:], "`")

		if matchStart == -1 {
			break
		}

		matchStart = matchStart + matchEnd
		matchEnd = matchStart + 1

		for matchEnd < max && CharCodeAt(state.Src, matchEnd) == 0x60 /* ` */ {
			matchEnd++
		}

		closerLength := matchEnd - matchStart

		if closerLength == openerLength {
			// Found matching closer length.
			if !silent {
				token := state.Push("code_inline", "code", 0)
				token.Markup = marker
				token.Content = state.Src[pos:matchStart]
				token.Content = NEWLINES_RE.ReplaceAllString(token.Content, " ")
				token.Content = BACKTICK_RE.ReplaceAllString(token.Content, "$1")
			}

			state.Pos = matchEnd
			return true
		}
		// Some different length found, put it in cache as upper limit of where closer can be found
		state.Backticks[closerLength] = matchStart
	}

	// Scanned through the end, didn't find anything
	state.BackTicksScanned = true

	if !silent {
		state.Pending += marker
	}

	state.Pos += openerLength

	return true
}
