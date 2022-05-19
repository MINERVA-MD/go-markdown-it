package pkg

func TextJoin(state *StateCore, _ *StateBlock, _ *StateInline, _ int, _ int, _ bool) bool {

	blockTokens := state.Tokens
	l := len(*blockTokens)
	for j := 0; j < l; j++ {
		if (*blockTokens)[j].Type != "inline" {
			continue
		}

		tokens := (*blockTokens)[j].Children
		max := len(tokens)

		for curr := 0; curr < max; curr++ {
			if (*blockTokens)[j].Children[curr].Type == "text_special" {
				(*blockTokens)[j].Children[curr].Type = "text"
			}
		}

		last := 0
		curr := 0
		for curr = 0; curr < max; curr++ {
			if tokens[curr].Type == "text" &&
				curr+1 < max &&
				tokens[curr+1].Type == "text" {

				// collapse two adjacent text nodes
				(*blockTokens)[j].Children[curr+1].Content = tokens[curr].Content + tokens[curr+1].Content
			} else {
				if curr != last {
					(*blockTokens)[j].Children[last] = tokens[curr]
				}

				last++
			}
		}

		if curr != last {
			tokens = tokens[:last]
			(*blockTokens)[j].Children = tokens
		}
	}

	return true
}
