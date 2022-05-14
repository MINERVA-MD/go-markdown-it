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
	//fmt.Println("Entered Backtick")

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

	marker := Slice(state.Src, start, pos)
	openerLength := utf8.RuneCountInString(marker)

	backTicksCheck := 0

	if len(state.Backticks) > 0 &&
		state.Backticks[openerLength] != 0 {
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

	//fmt.Println(start, ch, pos, marker, openerLength, matchStart, matchEnd)
	//fmt.Println(state.Src, matchEnd, strings.Index(state.Src, "`"))

	//fmt.Println(strings.Index(" b `", "`"))

	//fmt.Println(state.Src, matchEnd, strings.Index(state.Src, "`"))
	//fmt.Println(Slice(state.Src, matchEnd, utf8.RuneCountInString(state.Src)))
	//fmt.Println(IndexOfSubstring(Slice(state.Src, matchEnd, utf8.RuneCountInString(state.Src)), "`"))
	//fmt.Println(utf8.RuneCountInString(state.Src))
	//fmt.Println(Slice(state.Src, matchEnd, utf8.RuneCountInString(state.Src)))
	//fmt.Println(utf8.RuneCountInString(Slice(state.Src, matchEnd, utf8.RuneCountInString(state.Src))))

	//fmt.Println(state.Src, strings.Index(Slice(state.Src, matchEnd, utf8.RuneCountInString(state.Src)), "`"))
	// Nothing found in the cache, scan until the end of the line (or until marker is found)
	for {
		matchStart = IndexOfSubstring(Slice(state.Src, matchEnd, utf8.RuneCountInString(state.Src)), "`")
		if matchStart == -1 {
			break
		}

		//fmt.Println(matchStart, matchEnd)

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
				token.Content = Slice(state.Src, pos, matchStart)
				//fmt.Println(state.Src, pos, matchStart)
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

	//fmt.Println(start, ch, pos, marker, openerLength)

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
