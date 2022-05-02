package inline

import (
	. "go-markdown-it/pkg/common"
	. "go-markdown-it/pkg/maps"
	. "go-markdown-it/pkg/rules/block"
	"go-markdown-it/pkg/rules/core"
	"strconv"
	"strings"
)

func Entity(
	_ *core.StateCore,
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

	if CharCodeAt(state.Src, pos) != 0x26 /* & */ {
		return false
	}

	if pos+1 >= max {
		return false
	}

	ch := CharCodeAt(state.Src, pos+1)

	if ch == 0x23 /* # */ {
		match := DIGITAL_RE.FindStringSubmatch(state.Src[pos:])

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
			state.Pos += len(match[0])
			return true
		}
	} else {
		match := NAMED_RE.FindStringSubmatch(state.Src[pos:])

		if len(match) > 0 {
			if _, ok := ENTITIES[match[1]]; ok {
				if !silent {
					token := state.Push("text_special", "", 0)
					token.Content = ENTITIES[match[1]]
					token.Markup = match[0]
					token.Info = "entity"

					state.Pos += len(match[0])
					return true
				}
			}
		}
	}

	return false
}
