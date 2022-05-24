package pkg

import (
	"strings"
	"unicode/utf8"
)

func Backtick(
	_ *StateCore,
	_ *StateBlock,
	state *StateInline,
	_ int,
	_ int,
	silent bool,
) bool {
	return state.Backtick(silent)
}

func (state *StateInline) Backtick(silent bool) bool {

	pos := state.Pos
	ch, _ := state.Src2.CharCodeAt(pos)

	if ch != 0x60 /* ` */ {
		return false
	}

	start := pos
	pos++
	max := state.PosMax

	// scan marker length
	for {
		if cc, _ := state.Src2.CharCodeAt(pos); pos < max && cc == 0x60 /* ` */ {
			pos++
		} else {
			break
		}
	}

	marker := state.Src2.Slice(start, pos)
	openerLength := utf8.RuneCountInString(marker)

	backTicksCheck := 0

	if len(state.Backticks) > 0 &&
		state.Backticks[openerLength] != 0 {
		backTicksCheck = state.Backticks[openerLength]
	}

	if state.BackTicksScanned && (backTicksCheck <= start) {
		if !silent {
			_ = state.Pending2.WriteString(marker)
		}

		state.Pos += openerLength
		return true
	}

	matchStart := pos
	matchEnd := pos

	for {
		slice := state.Src2.Slice(matchEnd, state.Src2.Length)
		matchStart = IndexOfSubstring(slice, "`")
		if matchStart == -1 {
			break
		}

		matchStart = matchStart + matchEnd
		matchEnd = matchStart + 1

		for {
			if cc, _ := state.Src2.CharCodeAt(matchEnd); matchEnd < max && cc == 0x60 /* ` */ {
				matchEnd++
			} else {
				break
			}
		}

		closerLength := matchEnd - matchStart

		if closerLength == openerLength {
			// Found matching closer length.
			if !silent {
				token := state.Push("code_inline", "code", 0)
				token.Markup = marker
				token.Content = state.Src2.Slice(pos, matchStart)
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
		_ = state.Pending2.WriteString(marker)
	}

	state.Pos += openerLength

	return true
}

func IndexOfSubstring(s string, substr string) int {

	byteIndex := strings.Index(s, substr)
	if byteIndex < 0 {
		return -1
	}
	return utf8.RuneCountInString(s[:byteIndex])
}

func LastIndexOfSubstring(s string, substr string) int {

	byteIndex := strings.LastIndex(s, substr)
	if byteIndex < 0 {
		return -1
	}
	return utf8.RuneCountInString(s[:byteIndex])
}
