package pkg

import (
	"fmt"
	"strconv"
)

// Search `[-+*][\n ]`, returns next pos after marker on success
// or -1 on fail.

func (state *StateBlock) SkipBulletListMarker(startLine int) int {

	pos := state.BMarks[startLine] + state.TShift[startLine]
	max := state.EMarks[startLine]

	marker := CharCodeAt(state.Src, pos)
	pos++

	// Check bullet
	if marker != 0x2A /* * */ &&
		marker != 0x2D /* - */ &&
		marker != 0x2B {
		return -1
	}

	if pos < max {
		ch := CharCodeAt(state.Src, pos)

		if !IsSpace(ch) {
			// " -test " - is not a list item
			return -1
		}
	}

	return pos
}

func (state *StateBlock) SkipOrderedListMarker(startLine int) int {
	start := state.BMarks[startLine] + state.TShift[startLine]
	pos := start
	max := state.EMarks[startLine]

	// List marker should have at least 2 chars (digit + dot)
	if pos+1 >= max {
		return -1
	}

	ch := CharCodeAt(state.Src, pos)
	pos++

	if ch < 0x30 /* 0 */ || ch > 0x39 {
		return -1
	}

	for {
		// EOL -> fail
		if pos >= max {
			return -1
		}

		ch = CharCodeAt(state.Src, pos)
		pos++

		if ch >= 0x30 /* 0 */ && ch <= 0x39 {

			// List marker should have no more than 9 digits
			// (prevents integer overflow in browsers)
			if pos-start >= 10 {
				return -1
			}

			continue
		}

		// found valid marker
		if ch == 0x29 /* ) */ || ch == 0x2e /* . */ {
			break
		}

		return -1
	}

	if pos < max {
		ch = CharCodeAt(state.Src, pos)

		if !IsSpace(ch) {
			// " 1.test " - is not a list item
			return -1
		}
	}
	return pos
}

func (state *StateBlock) MarkTightParagraphs(idx int) {

	level := state.Level + 2
	l := len(*state.Tokens)

	for i := idx + 2; i < l; i++ {
		if (*state.Tokens)[i].Level == level && (*state.Tokens)[i].Type == "paragraph_open" {
			(*state.Tokens)[i+2].Hidden = true
			(*state.Tokens)[i].Hidden = true
			i += 2
		}
	}
}

func List(
	_ *StateCore,
	state *StateBlock,
	_ *StateInline,
	startLine int,
	endLine int,
	silent bool,
) bool {
	return state.List(startLine, endLine, silent)
}

func (state *StateBlock) List(startLine int, endLine int, silent bool) bool {

	fmt.Println("Processing List")
	var pos int
	var max int
	var start int
	tight := true
	var indent int
	var offset int
	var initial int
	var terminate bool
	var itemLines []int

	var isOrdered bool
	var markerValue int
	var contentStart int
	var indentAfterMarker int
	isTerminatingParagraph := false

	// if it"s indented more than 3 spaces, it should be a code block
	if state.SCount[startLine]-state.BlkIndent >= 4 {
		return false
	}

	// Special case:
	//  - item 1
	//   - item 2
	//    - item 3
	//     - item 4
	//      - this one is a paragraph continuation
	if state.ListIndent >= 0 &&
		state.SCount[startLine]-state.ListIndent >= 4 &&
		state.SCount[startLine] < state.BlkIndent {
		return false
	}

	// limit conditions when list can interrupt
	// a paragraph (validation mode only)
	if silent && state.ParentType == "paragraph" {
		// Next list item should still terminate previous list item;
		//
		// This code can fail if plugins use blkIndent as well as lists,
		// but I hope the spec gets fixed long before that happens.
		//
		if state.SCount[startLine] >= state.BlkIndent {
			isTerminatingParagraph = true
		}
	}

	// Detect list type and position after marker
	posAfterMarker := state.SkipOrderedListMarker(startLine)
	if posAfterMarker >= 0 {
		isOrdered = true
		start = state.BMarks[startLine] + state.TShift[startLine]
		markerValue, _ = strconv.Atoi(state.Src[start : posAfterMarker-1])

		// If we"re starting a new ordered list right after
		// a paragraph, it should start with 1.
		if isTerminatingParagraph && markerValue != 1 {
			return false
		}
	} else {
		posAfterMarker = state.SkipBulletListMarker(startLine)
		if posAfterMarker >= 0 {
			isOrdered = false

		} else {
			return false
		}
	}

	// If we"re starting a new unordered list right after
	// a paragraph, first line should not be empty.
	if isTerminatingParagraph {
		if state.SkipSpaces(posAfterMarker) >= state.EMarks[startLine] {
			return false
		}
	}

	// We should terminate list on style change. Remember first one to compare.
	markerCharCode := CharCodeAt(state.Src, posAfterMarker-1)

	// For validation mode we can terminate immediately
	if silent {
		return true
	}

	// Start list
	listTokIdx := len(*state.Tokens)
	var token *Token

	if isOrdered {
		token = state.Push("ordered_list_open", "ol", 1)
		if markerValue != 1 {
			token.Attrs = []Attribute{
				{
					Name: "start",
					// TODO: may need to check this
					Value: string(markerValue),
				},
			}
		}

	} else {
		token = state.Push("bullet_list_open", "ul", 1)
	}

	listLines := []int{startLine, 0}
	token.Map = []int{startLine, 0}
	token.Markup = string(markerCharCode)

	//
	// Iterate list items
	//

	nextLine := startLine
	prevEmptyEnd := false
	terminatorRules := state.Md.Block.Ruler.GetRules("list")

	oldParentType := state.ParentType
	state.ParentType = "list"

	for nextLine < endLine {
		pos = posAfterMarker
		max = state.EMarks[nextLine]

		offset = state.SCount[nextLine] + posAfterMarker - (state.BMarks[startLine] + state.TShift[startLine])
		initial = state.SCount[nextLine] + posAfterMarker - (state.BMarks[startLine] + state.TShift[startLine])

		for pos < max {
			ch := CharCodeAt(state.Src, pos)

			if ch == 0x09 {
				offset += 4 - (offset+state.BsCount[nextLine])%4
			} else if ch == 0x20 {
				offset++
			} else {
				break
			}

			pos++
		}

		contentStart = pos

		if contentStart >= max {
			// trimming space in "-    \n  3" case, indent is 1 here
			indentAfterMarker = 1
		} else {
			indentAfterMarker = offset - initial
		}

		// If we have more than 4 spaces, the indent is 1
		// (the rest is just indented code block)
		if indentAfterMarker > 4 {
			indentAfterMarker = 1
		}

		// "  -  test"
		//  ^^^^^ - calculating total length of this thing
		indent = initial + indentAfterMarker

		// Run subparser & write tokens
		token := state.Push("list_item_open", "li", 1)
		token.Markup = string(markerCharCode)

		itemLines = []int{startLine, 0}
		token.Map = []int{startLine, 0}
		if isOrdered {
			token.Info = state.Src[start : posAfterMarker-1]
		}

		// change current state, then restore it after parser subcall
		oldTight := state.Tight
		oldTShift := state.TShift[startLine]
		oldSCount := state.SCount[startLine]

		//  - example list
		// ^ listIndent position will be here
		//   ^ blkIndent position will be here
		//
		oldListIndent := state.ListIndent
		state.ListIndent = state.BlkIndent
		state.BlkIndent = indent

		state.Tight = true
		state.TShift[startLine] = contentStart - state.BMarks[startLine]
		state.SCount[startLine] = offset

		if contentStart >= max && state.IsEmpty(startLine+1) {
			// workaround for this case
			// (list item is empty, list terminates before "foo"):
			// ~~~~~~~~
			//   -
			//
			//     foo
			// ~~~~~~~~
			state.Line = Min(state.Line+2, endLine)
			//Math.min(state.Line+2, endLine)
		} else {
			state.Md.Block.Tokenize(state, startLine, endLine, true)
		}

		// If any of list item is tight, mark list as tight
		if !state.Tight || prevEmptyEnd {
			tight = false
		}
		// Item become loose if finish with empty line,
		// but we should filter last element, because it means list finish
		prevEmptyEnd = (state.Line-startLine) > 1 && state.IsEmpty(state.Line-1)

		state.BlkIndent = state.ListIndent
		state.ListIndent = oldListIndent
		state.TShift[startLine] = oldTShift
		state.SCount[startLine] = oldSCount
		state.Tight = oldTight

		token = state.Push("list_item_close", "li", -1)
		token.Markup = string(markerCharCode)

		startLine = state.Line
		nextLine = state.Line
		itemLines[1] = nextLine
		contentStart = state.BMarks[startLine]

		if nextLine >= endLine {
			break
		}

		//
		// Try to check if list is terminated or continued.
		//
		if state.SCount[nextLine] < state.BlkIndent {
			break
		}

		// if it's indented more than 3 spaces, it should be a code block
		if state.SCount[startLine]-state.BlkIndent >= 4 {
			break
		}

		// fail if terminating block found
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

		// fail if list has another type
		if isOrdered {
			posAfterMarker = state.SkipOrderedListMarker(nextLine)
			if posAfterMarker < 0 {
				break
			}
			start = state.BMarks[nextLine] + state.TShift[nextLine]
		} else {
			posAfterMarker = state.SkipBulletListMarker(nextLine)
			if posAfterMarker < 0 {
				break
			}
		}

		if markerCharCode != CharCodeAt(state.Src, posAfterMarker-1) {
			break
		}

	}

	// Finalize list
	if isOrdered {
		token = state.Push("ordered_list_close", "ol", -1)
	} else {
		token = state.Push("bullet_list_close", "ul", -1)
	}
	token.Markup = string(markerCharCode)

	listLines[1] = nextLine
	state.Line = nextLine

	state.ParentType = oldParentType

	// mark paragraphs tight if needed
	if tight {
		state.MarkTightParagraphs(listTokIdx)
	}

	return true
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
