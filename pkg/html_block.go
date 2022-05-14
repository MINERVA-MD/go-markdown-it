package pkg

import (
	"unicode/utf8"
)

func HtmlBlock(
	_ *StateCore,
	state *StateBlock,
	_ *StateInline,
	startLine int,
	endLine int,
	silent bool,
) bool {
	return state.HtmlBlock(startLine, endLine, silent)
}

func (state *StateBlock) HtmlBlock(startLine int, endLine int, silent bool) bool {
	//fmt.Println("Processing Html Block")

	pos := state.BMarks[startLine] + state.TShift[startLine]
	max := state.EMarks[startLine]

	//fmt.Println(pos, max, startLine, endLine)

	// if it's indented more than 3 spaces, it should be a code block
	if state.SCount[startLine]-state.BlkIndent >= 4 {
		return false
	}
	if !state.Md.Options.Html {
		return false
	}

	if CharCodeAt(state.Src, pos) != 0x3C {
		return false
	}

	lineText := Slice(state.Src, pos, max)

	var i = 0
	for _, sequence := range HTML_SEQUENCES {
		match, _ := sequence.Start.MatchString(lineText)
		if match {
			break
		}
		i++
	}

	if i == len(HTML_SEQUENCES) {
		return false
	}

	if silent {
		return HTML_SEQUENCES[i].Terminate
	}

	nextLine := startLine + 1

	matchEnd, _ := HTML_SEQUENCES[i].End.MatchString(lineText)
	if !matchEnd {
		for ; nextLine < endLine; nextLine++ {
			if state.SCount[nextLine] < state.BlkIndent {
				break
			}

			pos = state.BMarks[nextLine] + state.TShift[nextLine]
			max = state.EMarks[nextLine]
			lineText = Slice(state.Src, pos, max)

			matchEnd, _ = HTML_SEQUENCES[i].End.MatchString(lineText)
			if matchEnd {
				if utf8.RuneCountInString(lineText) != 0 {
					nextLine++
				}
				break
			}
		}
	}

	state.Line = nextLine

	token := state.Push("html_block", "", 0)
	token.Map = []int{startLine, nextLine}
	token.Content = state.GetLines(startLine, nextLine, state.BlkIndent, true)

	//fmt.Println(token.Map)
	return true
}
