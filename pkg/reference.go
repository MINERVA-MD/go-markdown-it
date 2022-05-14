package pkg

import (
	"strings"
	"unicode/utf8"
)

func Reference(
	_ *StateCore,
	state *StateBlock,
	_ *StateInline,
	startLine int,
	endLine int,
	silent bool,
) bool {
	return state.Reference(startLine, endLine, silent)
}

func (state *StateBlock) Reference(startLine int, _ int, silent bool) bool {

	//fmt.Println("Processing Reference")
	lines := 0
	var endLine int
	var labelEnd int
	var terminate bool
	pos := state.BMarks[startLine] + state.TShift[startLine]
	max := state.EMarks[startLine]
	nextLine := startLine + 1

	// if it's indented more than 3 spaces, it should be a code block
	if state.SCount[startLine]-state.BlkIndent >= 4 {
		return false
	}

	if CharCodeAt(state.Src, pos) != 0x5B /* [ */ {
		return false
	}

	// Simple check to quickly interrupt scan on [link](url) at the start of line.
	// Can be useful on practice: https://github.com/markdown-it/markdown-it/issues/54
	for {
		pos++
		if pos >= max {
			break
		}

		if CharCodeAt(state.Src, pos) == 0x5D /* ] */ &&
			CharCodeAt(state.Src, pos-1) != 0x5C /* \ */ {
			if pos+1 == max {
				return false
			}
			if CharCodeAt(state.Src, pos+1) != 0x3A {
				return false
			}
			break
		}
	}

	endLine = state.LineMax

	// jump line-by-line until empty one or EOF
	terminatorRules := state.Md.Block.Ruler.GetRules("reference")

	oldParentType := state.ParentType
	state.ParentType = "reference"

	for ; nextLine < endLine && !state.IsEmpty(nextLine); nextLine++ {
		// this would be a code block normally, but after paragraph
		// it's considered a lazy continuation regardless of what's there
		if state.SCount[nextLine]-state.BlkIndent > 3 {
			continue
		}

		// quirk for blockquotes, this line should already be checked by that rule
		if state.SCount[nextLine] < 0 {
			continue
		}

		// Some tags can terminate paragraph without empty line.
		terminate = false
		l := len(terminatorRules)
		for i := 0; i < l; i++ {
			if terminatorRules[i](nil, state, nil, nextLine, endLine, true) {
				terminate = true
				break
			}
		}
		if terminate {
			break
		}
	}

	str := strings.TrimSpace(state.GetLines(startLine, nextLine, state.BlkIndent, false))
	max = utf8.RuneCountInString(str)

	for pos = 1; pos < max; pos++ {
		ch := CharCodeAt(str, pos)
		if ch == 0x5B {
			return false
		} else if ch == 0x5D {
			labelEnd = pos
			break
		} else if ch == 0x0A {
			lines++
		} else if ch == 0x5C {
			pos++
			if pos < max && CharCodeAt(str, pos) == 0x0A {
				lines++
			}
		}
	}

	if labelEnd < 0 || CharCodeAt(str, labelEnd+1) != 0x3A {
		return false
	}

	// [label]:   destination   'title'
	//         ^^^ skip optional whitespace here
	for pos = labelEnd + 2; pos < max; pos++ {
		ch := CharCodeAt(str, pos)
		if ch == 0x0A {
			lines++
		} else if IsSpace(ch) {
			/*eslint no-empty:0*/
		} else {
			break
		}
	}

	// [label]:   destination   'title'
	//            ^^^^^^^^^^^ parse this
	res := state.Md.Helpers.ParseLinkDestination(str, pos, max)
	if !res.Ok {
		return false
	}

	href := state.Md.NormalizeLink(res.Str)
	if !state.Md.ValidateLink(href) {
		return false
	}

	//utils.PrettyPrint(href)

	pos = res.Pos
	lines += res.Lines

	// save cursor state, we could require to rollback later
	destEndPos := pos
	destEndLineNo := lines

	// [label]:   destination   'title'
	//                       ^^^ skipping those spaces
	start := pos
	for ; pos < max; pos++ {
		ch := CharCodeAt(str, pos)
		if ch == 0x0A {
			lines++
		} else if IsSpace(ch) {
			/*eslint no-empty:0*/
		} else {
			break
		}
	}

	// [label]:   destination   'title'
	//                          ^^^^^^^ parse this
	var title string
	res = state.Md.Helpers.ParseLinkTitle(str, pos, max)
	if pos < max && start != pos && res.Ok {
		title = res.Str
		pos = res.Pos
		lines += res.Lines
	} else {
		title = ""
		pos = destEndPos
		lines = destEndLineNo
	}

	// skip trailing spaces until the rest of the line
	for pos < max {
		ch := CharCodeAt(str, pos)
		if !IsSpace(ch) {
			break
		}
		pos++
	}

	if pos < max && CharCodeAt(str, pos) != 0x0A {
		if utf8.RuneCountInString(title) > 0 {
			// garbage at the end of the line after title,
			// but it could still be a valid reference if we roll back
			title = ""
			pos = destEndPos
			lines = destEndLineNo
			for pos < max {
				ch := CharCodeAt(str, pos)
				if !IsSpace(ch) {
					break
				}
				pos++
			}
		}
	}

	if pos < max && CharCodeAt(str, pos) != 0x0A {
		// garbage at the end of the line
		return false
	}

	label := NormalizeReference(Slice(str, 1, labelEnd))
	if utf8.RuneCountInString(label) == 0 {
		// CommonMark 0.20 disallows empty labels
		return false
	}

	// LinkReference can not terminate anything. This check is for safety only.
	if silent {
		return true
	}

	if state.Env.References == nil {
		state.Env.References = map[string]LinkReference{}
	}

	if _, ok := state.Env.References[label]; !ok {
		state.Env.References[label] = LinkReference{
			Href:  href,
			Title: title,
		}
	}

	state.ParentType = oldParentType

	state.Line = startLine + lines + 1
	return true
}
