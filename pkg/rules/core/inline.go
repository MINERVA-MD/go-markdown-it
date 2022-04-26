package core

import . "go-markdown-it/pkg/types"

func InlineCore(state *StateCore) {
	for _, token := range state.Tokens {
		if token.Type != "inline" {
			// TODO: state.md.inline.parse(tok.content, state.md, state.env, tok.children);
		}
	}
}
