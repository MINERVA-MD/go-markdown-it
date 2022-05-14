package pkg

import (
	"unicode/utf8"
)

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
	//fmt.Println("Processing Image")
	href := ""
	title := ""
	label := ""
	var code rune
	var start int
	oldPos := state.Pos
	max := state.PosMax
	var res LinkResult

	//fmt.Println("Entered Image")
	//fmt.Println(state.Src, state.Pos, CharCodeAt(state.Src, state.Pos))
	if CharCodeAt(state.Src, state.Pos) != 0x21 /* ! */ {
		//fmt.Println("Returning false 1")
		return false
	}
	if CharCodeAt(state.Src, state.Pos+1) != 0x5B /* [ */ {
		//fmt.Println("Returning false 2")
		return false
	}

	labelStart := state.Pos + 2
	labelEnd := state.Md.Helpers.ParseLinkLabel(state, state.Pos+1, false)

	// parser failed to find ']', so it's not a valid link
	if labelEnd < 0 {
		//fmt.Println("Returning false 3")
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
			//fmt.Println("Returning false 4")
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
			//fmt.Println("Returning false 5")
			return false
		}
		pos++
	} else {
		// Link reference
		if state.Env.References == nil {
			//fmt.Println("Returning false 6")
			return false
		}

		if pos < max && CharCodeAt(state.Src, pos) == 0x5B /* [ */ {
			start = pos + 1
			pos = state.Md.Helpers.ParseLinkLabel(state, pos, false)
			if pos >= 0 {
				label = Slice(state.Src, start, pos)
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
			label = Slice(state.Src, labelStart, labelEnd)
		}

		//utils.PrettyPrint(state.Env.References)
		normalizeReference := NormalizeReference(label)
		//fmt.Println(label)
		if _, ok := state.Env.References[normalizeReference]; !ok {
			state.Pos = oldPos
			//fmt.Println("Returning false 7")
			return false
		}

		ref := state.Env.References[normalizeReference]
		href = ref.Href
		title = ref.Title
	}

	// We found the end of the link, and know for a fact it's a valid link;
	// so all that's left to do is to call tokenizer.
	if !silent {

		content := Slice(state.Src, labelStart, labelEnd)

		var tokens []*Token
		state.Md.Inline.Parse(content, state.Md, state.Env, &tokens)

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

		if utf8.RuneCountInString(title) > 0 {
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
