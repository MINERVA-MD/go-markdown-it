package pkg

import "fmt"

func InlineCore(state *StateCore, _ *StateBlock, _ *StateInline, _ int, _ int, _ bool) bool {
	fmt.Println("Running Inline Core")
	for idx := 0; idx < len(state.Tokens); idx++ {
		if state.Tokens[idx].Type == "inline" {
			state.Md.Inline.Parse(state.Tokens[idx].Content, &state.Md, state.Env, &state.Tokens[idx].Children)
		}
	}
	return true
}
