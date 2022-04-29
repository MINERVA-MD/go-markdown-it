package core

import . "go-markdown-it/pkg/types"

func InlineCore(state *StateCore) {
	for _, token := range state.Tokens {
		if token.Type != "inline" {
			state.Md.Inline.Parse(token.Content, &state.Md, state.Env, &token.Children)
		}
	}
}
