package pkg

import (
	"unicode/utf8"
)

func IsLinkOpen(str string) bool {
	return LINK_OPEN.MatchString(str)
}

func IsLinkClose(str string) bool {
	return LINK_CLOSE.MatchString(str)
}

func IsLetter(ch rune) bool {
	var lc = ch | 0x20 // to lowercase
	return (lc >= 0x61 /* a */) && (lc <= 0x7a /* z */)
}

func HtmlInline(
	_ *StateCore,
	_ *StateBlock,
	state *StateInline,
	_ int,
	_ int,
	silent bool,
) bool {
	return state.HtmlInline(silent)
}

func (state *StateInline) HtmlInline(silent bool) bool {
	pos := state.Pos

	if !state.Md.Options.Html {
		return false
	}

	// Check Start
	max := state.PosMax
	if cc, _ := state.Src2.CharCodeAt(pos); cc != 0x3C /* < */ ||
		pos+2 >= max {
		return false
	}

	// Quick fail on second char
	ch, _ := state.Src2.CharCodeAt(pos + 1)
	if ch != 0x21 /* ! */ &&
		ch != 0x3F /* ? */ &&
		ch != 0x2F /* / */ &&
		!IsLetter(ch) {
		return false
	}

	// TODO: Replace wit Slice function
	slice := state.Src2.Slice(pos, state.Src2.Length)
	match := HTML_TAG_RE.FindStringSubmatch(slice)

	if len(match) == 0 {
		return false
	}

	matchLen := utf8.RuneCountInString(match[0])
	if !silent {
		token := state.Push("html_inline", "", 0)

		slice := state.Src2.Slice(pos, pos+matchLen)
		token.Content = slice

		if IsLinkOpen(token.Content) {
			state.LinkLevel++
		}

		if IsLinkClose(token.Content) {
			state.LinkLevel--
		}
	}

	state.Pos += matchLen

	return true
}
