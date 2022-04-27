package inline

import (
	. "go-markdown-it/pkg"
	. "go-markdown-it/pkg/common"
	. "go-markdown-it/pkg/rules/block"
	. "go-markdown-it/pkg/types"
)

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

	if CharCodeAt(state.Src, state.Pos) != 0x5B /* [ */ {
		return false
	}

	labelStart := state.Pos + 1
	labelEnd := state.Md.Helpers.ParseLinkLabel(state.Pos, true)

	// parser failed to find ']', so it's not a valid link
	if labelEnd < 0 {
		return false
	}

	pos := labelEnd + 1

	if pos < max && CharCodeAt(state.Src, pos) == 0x28 /* ( */ {
		// Inline link

		parseReference = false

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
			}
		}
		if pos >= max || CharCodeAt(state.Src, pos) != 0x29 /* ) */ {
			// parsing a valid shortcut link failed, fallback to reference
			parseReference = true
		}
		pos++
	}

	if parseReference {
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
		state.Pos = labelStart
		state.PosMax = labelEnd

		token := state.Push("link_open", "a", 1)
		token.Attrs = []Attribute{
			{
				Name:  "href",
				Value: href,
			},
		}

		if len(title) > 0 {
			token.Attrs = append(token.Attrs, Attribute{
				Name:  "title",
				Value: title,
			})
		}

		state.LinkLevel++
		state.Md.Inline.Tokenize(state)
		state.LinkLevel--

		token = state.Push("link_close", "a", -1)
	}

	state.Pos = pos
	state.PosMax = max

	return true
}
