package pkg

func TextJoin(state *StateCore, _ *StateBlock, _ *StateInline, _ int, _ int, _ bool) bool {

	blockTokens := state.Tokens
	l := len(blockTokens)
	for j := 0; j < l; j++ {
		if blockTokens[j].Type != "inline" {
			continue
		}

		tokens := blockTokens[j].Children
		max := len(tokens)

		for curr := 0; curr < max; curr++ {
			if tokens[curr].Type == "text_special" {
				tokens[curr].Type = "text"
			}
		}

		last := 0
		curr := 0
		for curr = 0; curr < max; curr++ {
			if tokens[curr].Type == "text" &&
				curr+1 < max &&
				tokens[curr+1].Type == "text" {

				// collapse two adjacent text nodes
				tokens[curr+1].Content = tokens[curr].Content + tokens[curr+1].Content
			} else {
				if curr != last {
					tokens[last] = tokens[curr]
				}

				last++
			}
		}

		if curr != last {
			tokens = tokens[:last]
		}
	}

	return true
}
