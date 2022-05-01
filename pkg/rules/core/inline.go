package core

import (
	. "go-markdown-it/pkg/rules/block"
	. "go-markdown-it/pkg/rules/inline"
)

func InlineCore(state *StateCore, _ *StateBlock, _ *StateInline, _ int, _ int, _ bool) bool {
	for _, token := range state.Tokens {
		if token.Type != "inline" {
			state.Md.Inline.Parse(token.Content, &state.Md, state.Env, &token.Children)
		}
	}
	return true
}
