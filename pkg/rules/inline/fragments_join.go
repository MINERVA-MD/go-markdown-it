package inline

// Clean up tokens after emphasis and strikethrough postprocessing:
// merge adjacent text nodes into one and re-calculate all token levels
//
// This is necessary because initially emphasis delimiter markers (*, _, ~)
// are treated as their own separate text tokens. Then emphasis rule either
// leaves them as text (needed to merge with adjacent text) or turns them
// into opening/closing tags (which messes up levels inside).
//

func (state *StateInline) FragmentsJoin() {
	curr := 0
	last := 0
	level := 0
	tokens := state.Tokens
	max := len(state.Tokens)

	for ; curr < max; curr++ {
		// re-calculate levels after emphasis/strikethrough turns some text nodes
		// into opening/closing tags
		token := tokens[curr]

		if token.Nesting < 0 { // closing tag
			level--
		}

		token.Level = level

		if token.Nesting > 0 {
			level++
		}

		if token.Type == "text" &&
			curr+1 < max &&
			tokens[curr+1].Type == "text" {
			// collapse two adjacent text nodes
			tokens[curr+1].Content = token.Content + tokens[curr+1].Content
		} else {
			if curr != last {
				tokens[last] = tokens[curr]
			}
			last++
		}
	}

	if curr != last {
		tokens = tokens[0:last]
	}
}
