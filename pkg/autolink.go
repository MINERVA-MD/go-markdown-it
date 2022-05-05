package pkg

func AutoLink(
	_ *StateCore,
	_ *StateBlock,
	state *StateInline,
	_ int,
	_ int,
	silent bool,
) bool {
	return state.AutoLink(silent)
}

func (state *StateInline) AutoLink(silent bool) bool {
	pos := state.Pos

	if CharCodeAt(state.Src, pos) != 0x3C {
		return false
	}

	start := state.Pos
	max := state.PosMax

	for {
		pos++
		if pos >= max {
			return false
		}

		ch := CharCodeAt(state.Src, pos)

		if ch == 0x3C /* < */ {
			return false
		}
		if ch == 0x3E /* > */ {
			break
		}
	}

	url := state.Src[start+1 : pos]

	if AUTOLINK_RE.MatchString(url) {
		fullUrl := state.Md.NormalizeLink(url)

		if !state.Md.ValidateLink(fullUrl) {
			return false
		}

		if !silent {
			token := state.Push("link_open", "a", 1)
			token.Attrs = []Attribute{
				{
					Name:  "href",
					Value: fullUrl,
				},
			}
			token.Markup = "autolink"
			token.Info = "auto"

			token = state.Push("text", "", 0)
			token.Content = state.Md.NormalizeLinkText(url)

			token = state.Push("link_close", "a", -1)
			token.Markup = "autolink"
			token.Info = "auto"
		}

		state.Pos += len(url) + 2
		return true
	}

	if EMAIL_RE.MatchString(url) {
		fullUrl := state.Md.NormalizeLink("mailto:" + url)
		if !state.Md.ValidateLink(fullUrl) {
			return false
		}

		if !silent {
			token := state.Push("link_open", "a", 1)
			token.Attrs = []Attribute{
				{
					Name:  "href",
					Value: fullUrl,
				},
			}
			token.Markup = "autolink"
			token.Info = "auto"

			token = state.Push("text", "", 0)
			token.Content = state.Md.NormalizeLinkText(url)

			token = state.Push("link_close", "a", -1)
			token.Markup = "autolink"
			token.Info = "auto"
		}

		state.Pos += len(url) + 2
		return true
	}

	return false
}
