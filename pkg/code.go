package pkg

import "fmt"

func Code(
	_ *StateCore,
	state *StateBlock,
	_ *StateInline,
	startLine int,
	endLine int,
	_ bool,
) bool {
	return state.Code(startLine, endLine)
}

func (state *StateBlock) Code(startLine int, endLine int) bool {

	fmt.Println("Processing Code")

	if state.SCount[startLine]-state.BlkIndent < 4 {
		return false
	}

	last := startLine + 1
	nextLine := last

	for nextLine < endLine {
		fmt.Println(nextLine)
		fmt.Println(endLine)
		if state.IsEmpty(nextLine) {
			nextLine++
			continue
		}

		if state.SCount[nextLine]-state.BlkIndent >= 4 {
			nextLine++
			last = nextLine
			continue
		}
		break
	}

	state.Line = last

	token := state.Push("code_block", "code", 0)
	token.Content = state.GetLines(startLine, last, 4+state.BlkIndent, false) + "\n"
	token.Map = []int{startLine, state.Line}

	return true
}
