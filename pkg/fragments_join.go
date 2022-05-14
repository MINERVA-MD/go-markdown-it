package pkg

func FragmentsJoin(
	_ *StateCore,
	_ *StateBlock,
	state *StateInline,
	_ int,
	_ int,
	_ bool,
) bool {
	//fmt.Println("Processing FragmentsJoin")
	state.FragmentsJoin()
	return true
}

// Clean up tokens after emphasis and strikethrough postprocessing:
// merge adjacent text nodes into one and re-calculate all token levels
//
// This is necessary because initially emphasis delimiter markers (*, _, ~)
// are treated as their own separate text tokens. Then emphasis rule either
// leaves them as text (needed to merge with adjacent text) or turns them
// into opening/closing tags (which messes up levels inside).
//

func (state *StateInline) FragmentsJoin() {
	//fmt.Println("Running FragmentsJoin")
	curr := 0
	last := 0
	level := 0
	tokens := state.Tokens
	max := len(*state.Tokens)

	//fmt.Println(len(*state.Tokens))
	//utils.PrettyPrint(*state.Tokens)
	for ; curr < max; curr++ {
		// re-calculate levels after emphasis/strikethrough turns some text nodes
		// into opening/closing tags
		token := (*tokens)[curr]

		if token.Nesting < 0 { // closing tag
			level--
		}

		token.Level = level

		if token.Nesting > 0 {
			level++
		}

		if token.Type == "text" &&
			curr+1 < max &&
			(*state.Tokens)[curr+1].Type == "text" {
			// collapse two adjacent text nodes
			(*state.Tokens)[curr+1].Content = token.Content + (*tokens)[curr+1].Content
		} else {
			if curr != last {
				(*state.Tokens)[last] = (*tokens)[curr]
			}
			last++
		}
	}

	if curr != last {
		_tokens := (*tokens)[0:last]
		tokens = &_tokens
		*state.Tokens = _tokens
	}
	//fmt.Println(len(*state.Tokens))
	//utils.PrettyPrint(state.Tokens)
}
