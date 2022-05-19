package pkg

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

func Entity(
	_ *StateCore,
	_ *StateBlock,
	state *StateInline,
	_ int,
	_ int,
	silent bool,
) bool {
	return state.Entity(silent)
}

func (state *StateInline) Entity(silent bool) bool {

	pos := state.Pos
	max := state.PosMax

	if cc, _ := state.Src2.CharCodeAt(pos); cc != 0x26 /* & */ {
		return false
	}

	if pos+1 >= max {
		return false
	}

	ch, _ := state.Src2.CharCodeAt(pos + 1)

	if ch == 0x23 /* # */ {
		slice := state.Src2.Slice(pos, state.Src2.Length)
		match := DIGITAL_RE.FindStringSubmatch(slice)

		if len(match) > 0 {
			if !silent {
				firstChar := strings.ToLower(string(match[1][0]))

				var code rune
				if strings.ToLower(firstChar) == "x" {
					integer, _ := strconv.ParseInt(match[1][1:], 16, 0)
					code = rune(integer)
				} else {
					integer, _ := strconv.ParseInt(match[1], 10, 0)
					code = rune(integer)
				}

				token := state.Push("text_special", "", 0)

				if IsValidEntityCode(code) {
					token.Content = FromCodePoint(code)
				} else {
					token.Content = FromCodePoint(0xFFFD)
				}

				token.Markup = match[0]
				token.Info = "entity"
			}
			state.Pos += utf8.RuneCountInString(match[0])
			return true
		}
	} else {
		// TODO: Replace wit Slice function
		slice := state.Src2.Slice(pos, state.Src2.Length)
		match := NAMED_RE.FindStringSubmatch(slice)

		if len(match) > 0 {
			if _, ok := ENTITIES[match[1]]; ok {
				if !silent {
					token := state.Push("text_special", "", 0)
					token.Content = ENTITIES[match[1]]
					token.Markup = match[0]
					token.Info = "entity"

					state.Pos += utf8.RuneCountInString(match[0])
					return true
				}
			}
		}
	}

	return false
}
