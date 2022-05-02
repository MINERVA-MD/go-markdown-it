package block

import (
	"go-markdown-it/pkg"
	"go-markdown-it/pkg/rules/core"
	"go-markdown-it/pkg/rules/inline"
	"strings"
)

func Paragraph(
	_ *core.StateCore,
	state *StateBlock,
	_ *inline.StateInline,
	startLine int,
	_ int,
	_ bool,
) bool {
	return state.Paragraph(startLine)
}

func (state *StateBlock) Paragraph(startLine int) bool {
	nextLine := startLine + 1
	terminatorRules := state.Md.Block.Ruler.GetRules("paragraph")
	endLine := state.LineMax
	oldParentType := state.ParentType
	state.ParentType = "paragraph"

	var terminate bool

	// jump line-by-line until empty one or EOF
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

	content := strings.TrimSpace(state.GetLines(startLine, nextLine, state.BlkIndent, false))

	state.Line = nextLine

	token := state.Push("paragraph_open", "p", 1)
	token.Map = []int{startLine, state.Line}

	token = state.Push("inline", "", 0)
	token.Content = content
	token.Map = []int{startLine, state.Line}
	token.Children = []*pkg.Token{}

	token = state.Push("paragraph_close", "p", -1)

	state.ParentType = oldParentType

	return true
}
