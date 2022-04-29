package block

import (
	. "go-markdown-it/pkg"
	"strings"
)

func (state *StateBlock) LHeading(startLine int, endLine int) bool {

	var pos int
	var max int
	var level int
	var marker rune
	var terminate bool
	var content string
	var oldParentType string
	nextLine := startLine + 1
	terminatorRules := state.Md.Block.Ruler.GetRules("paragraph")

	// if it's indented more than 3 spaces, it should be a code block
	if state.SCount[startLine]-state.BlkIndent >= 4 {
		return false
	}

	oldParentType = state.ParentType
	state.ParentType = "paragraph" // use paragraph to match terminatorRules'

	// jump line-by-line until empty one or EOF
	for ; nextLine < endLine && !state.IsEmpty(nextLine); nextLine++ {
		// this would be a code block normally, but after paragraph
		// it's considered a lazy continuation regardless of what's there
		if state.SCount[nextLine]-state.BlkIndent > 3 {
			continue
		}

		//
		// Check for underline in setext header
		//
		if state.SCount[nextLine] >= state.BlkIndent {
			pos = state.BMarks[nextLine] + state.TShift[nextLine]
			max = state.EMarks[nextLine]

			if pos < max {
				marker = CharCodeAt(state.Src, pos)

				if marker == 0x2D /* - */ || marker == 0x3D /* = */ {
					pos = state.SkipChars(pos, marker)
					pos = state.SkipSpaces(pos)

					if pos >= max {
						if marker == 0x3D /* = */ {
							level = 1
						} else {
							level = 2
						}
						break
					}
				}
			}
		}

		// quirk for blockquotes, this line should already be checked by that rule
		if state.SCount[nextLine] < 0 {
			continue
		}

		// Some tags can terminate paragraph without empty line.
		terminate = false
		l := len(terminatorRules)
		for i := 0; i < l; i++ {
			if terminatorRules[i](state, nextLine, endLine, true) {
				terminate = true
				break
			}
		}
		if terminate {
			break
		}
	}

	if level == 0 {
		// Didn't find valid underline
		return false
	}

	content = strings.TrimSpace(state.GetLines(startLine, nextLine, state.BlkIndent, false))

	state.Line = nextLine + 1

	token := state.Push("heading_open", "h"+string(rune(level)), 1)
	token.Markup = string(marker)
	token.Map = []int{startLine, state.Line}

	token = state.Push("inline", "", 0)
	token.Content = content
	token.Map = []int{startLine, state.Line - 1}
	token.Children = []*Token{}

	token = state.Push("heading_close", "h"+string(rune(level)), -1)
	token.Markup = string(marker)

	state.ParentType = oldParentType

	return true
}
