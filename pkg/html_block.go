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
	pos := state.BMarks[startLine] + state.TShift[startLine]
	max := state.EMarks[startLine]

	// if it's indented more than 3 spaces, it should be a code block
	if state.SCount[startLine]-state.BlkIndent >= 4 {
		return false
	}
	if !state.Md.Options.Html {
		return false
	}

	if cc, _ := state.Src2.CharCodeAt(pos); cc != 0x3C {
		return false
	}

	lineText := state.Src2.Slice(pos, max)
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
			lineText = state.Src2.Slice(pos, max)

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

	return true
}
