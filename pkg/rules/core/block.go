package core

import (
	. "go-markdown-it/pkg"
	. "go-markdown-it/pkg/types"
)

func BlockCore(state *StateCore) {
	if state.InlineMode {
		var token = GenerateToken("inline", "", 0)
		token.Content = state.Src
		token.Map = []int{0, 1}
		token.Children = []*Token{}

		state.Tokens = append(state.Tokens, &token)
	} else {
		state.Md.Block.Parse(state.Src, &state.Md, state.Env, &state.Tokens)
	}
}
