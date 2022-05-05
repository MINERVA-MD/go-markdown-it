package pkg

func BlockCore(state *StateCore, _ *StateBlock, _ *StateInline, _ int, _ int, _ bool) bool {
	if state.InlineMode {
		var token = GenerateToken("inline", "", 0)
		token.Content = state.Src
		token.Map = []int{0, 1}
		token.Children = []*Token{}
		state.Tokens = append(state.Tokens, &token)
	} else {
		state.Md.Block.Parse(state.Src, &state.Md, state.Env, &state.Tokens)
	}
	return true
}
