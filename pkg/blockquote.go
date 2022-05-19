package pkg

func BlockQuote(
	_ *StateCore,
	state *StateBlock,
	_ *StateInline,
	startLine int,
	endLine int,
	silent bool,
) bool {
	return state.BlockQuote(startLine, endLine, silent)
}

func (state *StateBlock) BlockQuote(startLine int, endLine int, silent bool) bool {

	var nextLine int
	var oldIndent int
	var adjustTab bool
	var terminate bool
	var oldBMarks []int
	var oldTShift []int
	var oldBSCount []int
	var spaceAfterMarker bool
	oldLineMax := state.LineMax
	pos := state.BMarks[startLine] + state.TShift[startLine]
	max := state.EMarks[startLine]

	// if it's indented more than 3 spaces, it should be a code block
	if state.SCount[startLine]-state.BlkIndent >= 4 {
		return false
	}

	// check the block quote marker
	if cc, _ := state.Src2.CharCodeAt(pos); cc != 0x3E /* > */ {
		return false
	}

	pos++

	// we know that it"s going to be a valid blockquote,
	// so no point trying to find the end of it in silent mode
	if silent {
		return true
	}

	// set offset past spaces and ">"
	initial := state.SCount[startLine] + 1
	offset := state.SCount[startLine] + 1

	// skip one optional space after ">"
	if cc, _ := state.Src2.CharCodeAt(pos); cc == 0x20 /* space */ {
		// " >   test "
		//     ^ -- position start of line here:
		pos++
		initial++
		offset++
		adjustTab = false
		spaceAfterMarker = true
	} else if cc, _ := state.Src2.CharCodeAt(pos); cc == 0x09 /* tab */ {
		spaceAfterMarker = true

		if (state.BsCount[startLine]+offset)%4 == 3 {
			// "  >\t  test "
			//       ^ -- position start of line here (tab has width===1)
			pos++
			initial++
			offset++
			adjustTab = false
		} else {
			// " >\t  test "
			//    ^ -- position start of line here + shift bsCount slightly
			//         to make extra space appear
			adjustTab = true
		}
	} else {
		spaceAfterMarker = false
	}

	oldBMarks = []int{state.BMarks[startLine]}
	state.BMarks[startLine] = pos

	for pos < max {
		ch, _ := state.Src2.CharCodeAt(pos)

		if IsSpace(ch) {
			if ch == 0x09 {
				adjustTabOffset := 0
				if adjustTab {
					adjustTabOffset = 1
				}
				offset += 4 - (offset+state.BsCount[startLine]+(adjustTabOffset))%4
			} else {
				offset++
			}
		} else {
			break
		}
		pos++
	}

	oldBSCount = []int{state.BsCount[startLine]}
	state.BsCount[startLine] = state.SCount[startLine] + 1 + 0

	if spaceAfterMarker {
		state.BsCount[startLine] += 1
	}

	lastLineEmpty := pos >= max
	oldSCount := []int{state.SCount[startLine]}

	state.SCount[startLine] = offset - initial
	oldTShift = []int{state.TShift[startLine]}

	state.TShift[startLine] = pos - state.BMarks[startLine]

	terminatorRules := state.Md.Block.Ruler.GetRules("blockquote")

	oldParentType := state.ParentType
	state.ParentType = "blockquote"

	// Search the end of the block
	//
	// Block ends with either:
	//  1. an empty line outside:
	//     ```
	//     > test
	//
	//     ```
	//  2. an empty line inside:
	//     ```
	//     >
	//     test
	//     ```
	//  3. another tag:
	//     ```
	//     > test
	//      - - -
	//     ```
	for nextLine = startLine + 1; nextLine < endLine; nextLine++ {

		// check if it's outdented, i.e. it's inside list item and indented
		// less than said list item:
		//
		// ```
		// 1. anything
		//    > current blockquote
		// 2. checking this line
		// ```
		isOutdented := state.SCount[nextLine] < state.BlkIndent

		pos = state.BMarks[nextLine] + state.TShift[nextLine]
		max = state.EMarks[nextLine]

		if pos >= max {
			// Case 1: line is not inside the blockquote, and this line is empty.
			break
		}

		isAlreadyIncremented := false
		if cc, _ := state.Src2.CharCodeAt(pos); cc == 0x3E /* > */ && !isOutdented {
			// This line is inside the blockquote.
			pos++
			isAlreadyIncremented = true

			// set offset past spaces and ">"
			offset = state.SCount[nextLine] + 1
			initial = state.SCount[nextLine] + 1

			// skip one optional space after '>'
			if cc, _ := state.Src2.CharCodeAt(pos); cc == 0x20 /* space */ {
				// ' >   test '
				//     ^ -- position start of line here:
				pos++
				initial++
				offset++
				adjustTab = false
				spaceAfterMarker = true
			} else if cc, _ := state.Src2.CharCodeAt(pos); cc == 0x09 /* tab */ {
				spaceAfterMarker = true

				if (state.BsCount[nextLine]+offset)%4 == 3 {
					// '  >\t  test '
					//       ^ -- position start of line here (tab has width===1)
					pos++
					initial++
					offset++
					adjustTab = false
				} else {
					// ' >\t  test '
					//    ^ -- position start of line here + shift bsCount slightly
					//         to make extra space appear
					adjustTab = true
				}
			} else {
				spaceAfterMarker = false
			}

			oldBMarks = append(oldBMarks, state.BMarks[nextLine])
			state.BMarks[nextLine] = pos

			for pos < max {
				ch, _ := state.Src2.CharCodeAt(pos)

				if IsSpace(ch) {
					if ch == 0x09 {
						adjustTabOffset := 0
						if adjustTab {
							adjustTabOffset = 1
						}

						offset += 4 - (offset+state.BsCount[nextLine]+(adjustTabOffset))%4
					} else {
						offset++
					}
				} else {
					break
				}

				pos++
			}

			lastLineEmpty = pos >= max

			oldBSCount = append(oldBSCount, state.BsCount[nextLine])
			state.BsCount[nextLine] = state.SCount[nextLine] + 1

			if spaceAfterMarker {
				state.BsCount[nextLine] += 1
			}

			oldSCount = append(oldSCount, state.SCount[nextLine])
			state.SCount[nextLine] = offset - initial

			oldTShift = append(oldTShift, state.TShift[nextLine])
			state.TShift[nextLine] = pos - state.BMarks[nextLine]
			continue
		}

		if !isAlreadyIncremented {
			pos++
		}

		// Case 2: line is not inside the blockquote, and the last line was empty.
		if lastLineEmpty {
			break
		}

		// Case 3: another tag found.
		terminate = false
		l := len(terminatorRules)
		for i := 0; i < l; i++ {
			if terminatorRules[i](nil, state, nil, nextLine, endLine, true) {
				terminate = true
				break
			}
		}

		if terminate {
			// Quirk to enforce "hard termination mode" for paragraphs;
			// normally if you call `tokenize(state, startLine, nextLine)`,
			// paragraphs will look below nextLine for paragraph continuation,
			// but if blockquote is terminated by another tag, they shouldn't
			state.LineMax = nextLine

			if state.BlkIndent != 0 {
				// state.blkIndent was non-zero, we now set it to zero,
				// so we need to re-calculate all offsets to appear as
				// if indent wasn't changed
				oldBMarks = append(oldBMarks, state.BMarks[nextLine])
				oldBSCount = append(oldBSCount, state.BsCount[nextLine])
				oldTShift = append(oldTShift, state.TShift[nextLine])
				oldSCount = append(oldSCount, state.SCount[nextLine])
				state.SCount[nextLine] -= state.BlkIndent
			}

			break
		}

		oldBMarks = append(oldBMarks, state.BMarks[nextLine])
		oldBSCount = append(oldBSCount, state.BsCount[nextLine])
		oldTShift = append(oldTShift, state.TShift[nextLine])
		oldSCount = append(oldSCount, state.SCount[nextLine])

		// A negative indentation means that this is a paragraph continuation
		//
		state.SCount[nextLine] = -1
	}

	oldIndent = state.BlkIndent
	state.BlkIndent = 0

	token := state.Push("blockquote_open", "blockquote", 1)
	token.Markup = ">"

	lines := []int{startLine, 0}
	token.Map = []int{startLine, 0}

	state.Md.Block.Tokenize(state, startLine, nextLine, false)

	token = state.Push("blockquote_close", "blockquote", -1)
	token.Markup = ">"

	state.LineMax = oldLineMax
	state.ParentType = oldParentType
	lines[1] = state.Line

	// Restore original tShift this might not be necessary since the parser
	// has already been here, but just to make sure we can do that.
	for i := 0; i < len(oldTShift); i++ {
		state.BMarks[i+startLine] = oldBMarks[i]
		state.TShift[i+startLine] = oldTShift[i]
		state.SCount[i+startLine] = oldSCount[i]
		state.BsCount[i+startLine] = oldBSCount[i]
	}

	state.BlkIndent = oldIndent

	return true
}
