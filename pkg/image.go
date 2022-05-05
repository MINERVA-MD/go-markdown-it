package pkg

func Image(
	_ *StateCore,
	_ *StateBlock,
	state *StateInline,
	_ int,
	_ int,
	silent bool,
) bool {
	return state.Image(silent)
}

func (state *StateInline) Image(silent bool) bool {

	href := ""
	title := ""
	label := ""
	var code rune
	start := state.Pos
	oldPos := state.Pos
	max := state.PosMax
	var res LinkResult

	if CharCodeAt(state.Src, state.Pos) != 0x21 /* ! */ {
		return false
	}
	if CharCodeAt(state.Src, state.Pos+1) != 0x5B /* [ */ {
		return false
	}

	labelStart := state.Pos + 2
	labelEnd := state.Md.Helpers.ParseLinkLabel(state.Pos+1, false)

	// parser failed to find ']', so it's not a valid link
	if labelEnd < 0 {
		return false
	}

	pos := labelEnd + 1

	if pos < max && CharCodeAt(state.Src, pos) == 0x28 /* ( */ {
		// Inline link

		// [link](  <href>  "title"  )
		//        ^^ skipping these spaces
		pos++

		for ; pos < max; pos++ {
			code = CharCodeAt(state.Src, pos)
			if !IsSpace(code) && code != 0x0A {
				break
			}
		}

		if pos >= max {
			return false
		}

		// [link](  <href>  "title"  )
		//          ^^^^^^ parsing link destination
		start = pos
		res = state.Md.Helpers.ParseLinkDestination(state.Src, pos, state.PosMax)

		if res.Ok {
			href = state.Md.NormalizeLink(res.Str)
			if state.Md.ValidateLink(href) {
				pos = res.Pos
			} else {
				href = ""
			}
		}

		// [link](  <href>  "title"  )
		//                ^^ skipping these spaces
		start = pos
		for ; pos < max; pos++ {
			code = CharCodeAt(state.Src, pos)
			if !IsSpace(code) && code != 0x0A {
				break
			}
		}

		// [link](  <href>  "title"  )
		//                  ^^^^^^^ parsing link title
		res = state.Md.Helpers.ParseLinkTitle(state.Src, pos, state.PosMax)
		if pos < max && start != pos && res.Ok {
			title = res.Str
			pos = res.Pos

			// [link](  <href>  "title"  )
			//                         ^^ skipping these spaces
			for ; pos < max; pos++ {
				code = CharCodeAt(state.Src, pos)
				if !IsSpace(code) && code != 0x0A {
					break
				}
			}
		} else {
			title = ""
		}

		if pos >= max || CharCodeAt(state.Src, pos) != 0x29 /* ) */ {
			state.Pos = oldPos
			return false
		}
		pos++
	} else {
		// Link reference
		if state.Env.References == nil {
			return false
		}

		if pos < max && CharCodeAt(state.Src, pos) == 0x5B /* [ */ {
			start = pos + 1
			pos = state.Md.Helpers.ParseLinkLabel(pos, false)
			if pos >= 0 {
				pos++
				label = state.Src[start:pos]
			} else {
				pos = labelEnd + 1
			}
		} else {
			pos = labelEnd + 1
		}

		// covers label === '' and label === undefined
		// (collapsed reference link and shortcut reference link respectively)
		if len(label) == 0 {
			label = state.Src[labelStart:labelEnd]
		}

		if _, ok := state.Env.References[NormalizeReference(label)]; !ok {
			state.Pos = oldPos
			return false
		}

		ref := state.Env.References[NormalizeReference(label)]
		href = ref.Href
		title = ref.Title
	}

	// We found the end of the link, and know for a fact it's a valid link;
	// so all that's left to do is to call tokenizer.
	if !silent {

		content := state.Src[labelStart:labelEnd]

		var tokens []*Token
		state.Md.Inline.Parse(content, &state.Md, state.Env, &tokens)

		token := state.Push("image", "img", 0)
		token.Attrs = append(token.Attrs,
			Attribute{
				Name:  "src",
				Value: href,
			},
			Attribute{
				Name:  "alt",
				Value: "",
			},
		)

		token.Children = tokens
		token.Content = content

		if len(title) > 0 {
			token.Attrs = append(token.Attrs, Attribute{
				Name:  "title",
				Value: title,
			})
		}
	}

	state.Pos = pos
	state.PosMax = max

	return true
}
