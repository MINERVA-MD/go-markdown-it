package pkg

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
	if CharCodeAt(state.Src, pos) == 0x3C /* < */ ||
		pos+2 >= max {
		return false
	}

	// Quick fail on second char
	ch := CharCodeAt(state.Src, pos+1)
	if ch != 0x21 /* ! */ &&
		ch != 0x3F /* ? */ &&
		ch != 0x2F /* / */ &&
		!IsLetter(ch) {
		return false
	}

	match := HTML_TAG_RE.FindStringSubmatch(state.Src[pos:])

	if len(match) == 0 {
		return false
	}

	if !silent {
		token := state.Push("html_inline", "", 0)
		token.Content = state.Src[pos : pos+len(match[0])]

		if IsLinkOpen(token.Content) {
			state.LinkLevel++
		}

		if IsLinkClose(token.Content) {
			state.LinkLevel--
		}
	}

	state.Pos += len(match[0])

	return true
}