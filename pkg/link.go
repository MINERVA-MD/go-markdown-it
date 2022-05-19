package pkg

import (
	"unicode/utf8"
)

func Link(
	_ *StateCore,
	_ *StateBlock,
	state *StateInline,
	_ int,
	_ int,
	silent bool,
) bool {
	return state.Link(silent)
}

func (state *StateInline) Link(silent bool) bool {

	href := ""
	title := ""
	label := ""
	var code rune
	start := state.Pos
	oldPos := state.Pos
	max := state.PosMax
	parseReference := true
	var res LinkResult

	//fmt.Println("Entered Link")

	if cc, _ := state.Src2.CharCodeAt(state.Pos); cc != 0x5B /* [ */ {
		return false
	}

	labelStart := state.Pos + 1
	labelEnd := state.Md.Helpers.ParseLinkLabel(state, state.Pos, true)

	// parser failed to find ']', so it's not a valid link
	if labelEnd < 0 {
		return false
	}

	pos := labelEnd + 1

	if cc, _ := state.Src2.CharCodeAt(pos); pos < max && cc == 0x28 /* ( */ {
		// Inline link

		parseReference = false

		// [link](  <href>  "title"  )
		//        ^^ skipping these spaces
		pos++

		for ; pos < max; pos++ {
			code, _ = state.Src2.CharCodeAt(pos)
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
		res = state.Md.Helpers.ParseLinkDestination(state.Src2, pos, state.PosMax)

		//fmt.Println(state.Src, pos, state.PosMax)
		//utils.PrettyPrint(res)
		if res.Ok {
			href = state.Md.NormalizeLink(res.Str)

			if state.Md.ValidateLink(href) {
				pos = res.Pos
			} else {
				href = ""
			}

			// [link](  <href>  "title"  )
			//                ^^ skipping these spaces
			start = pos
			for ; pos < max; pos++ {
				code, _ = state.Src2.CharCodeAt(pos)
				if !IsSpace(code) && code != 0x0A {
					break
				}
			}

			// [link](  <href>  "title"  )
			//                  ^^^^^^^ parsing link title
			res = state.Md.Helpers.ParseLinkTitle(state.Src2, pos, state.PosMax)
			if pos < max && start != pos && res.Ok {
				title = res.Str
				pos = res.Pos

				// [link](  <href>  "title"  )
				//                         ^^ skipping these spaces
				for ; pos < max; pos++ {
					code, _ = state.Src2.CharCodeAt(pos)
					if !IsSpace(code) && code != 0x0A {
						break
					}
				}
			}
		}

		if cc, _ := state.Src2.CharCodeAt(pos); pos >= max || cc != 0x29 /* ) */ {
			// parsing a valid shortcut link failed, fallback to reference
			parseReference = true
		}
		pos++
	}

	if parseReference {
		// Link reference

		//utils.PrettyPrint(state.Env.References)
		if state.Env.References == nil {
			return false
		}

		if cc, _ := state.Src2.CharCodeAt(pos); pos < max && cc == 0x5B /* [ */ {
			start = pos + 1
			pos = state.Md.Helpers.ParseLinkLabel(state, pos, false)
			if pos >= 0 {
				label = state.Src2.Slice(start, pos)
				pos++
			} else {
				pos = labelEnd + 1
			}
		} else {
			pos = labelEnd + 1
		}

		// covers label === '' and label === undefined
		// (collapsed reference link and shortcut reference link respectively)
		if utf8.RuneCountInString(label) == 0 {
			label = state.Src2.Slice(labelStart, labelEnd)
		}

		//fmt.Println(label, labelStart, labelEnd, NormalizeReference(label))
		//utils.PrettyPrint(state.Env.References)

		//str = strings.Replace(str, "", "", -1)

		//fmt.Println("1")
		//fmt.Println("label", label, strings.ToLower(label), strings.ToUpper(strings.ToLower(label)))
		// TODO: Refactor this into NormalizeReference
		normalizedReference := NormalizeReference(label)
		//fmt.Println("Normalized Reference", normalizedReference)
		if _, ok := state.Env.References[normalizedReference]; !ok {
			state.Pos = oldPos
			return false
		}

		//fmt.Println(label)
		//fmt.Println("2")
		//fmt.Println("label", label)
		ref := state.Env.References[normalizedReference]
		href = ref.Href
		title = ref.Title
	}

	// We found the end of the link, and know for a fact it's a valid link;
	// so all that's left to do is to call tokenizer.
	if !silent {
		state.Pos = labelStart
		state.PosMax = labelEnd

		token := state.Push("link_open", "a", 1)
		token.Attrs = []Attribute{
			{
				Name:  "href",
				Value: href,
			},
		}

		if utf8.RuneCountInString(title) > 0 {
			token.Attrs = append(token.Attrs, Attribute{
				Name:  "title",
				Value: title,
			})
		}

		state.LinkLevel++
		state.Md.Inline.Tokenize(state)
		state.LinkLevel--

		token = state.Push("link_close", "a", -1)

		//utils.PrettyPrint(state.Delimiters)
	}

	state.Pos = pos
	state.PosMax = max

	return true
}
