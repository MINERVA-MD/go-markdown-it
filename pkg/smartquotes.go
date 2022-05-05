package pkg

func ReplaceAt(str string, index int, ch rune) string {
	return str[0:index] + string(ch) + str[index+1:]
}

func ProcessInline(tokens []*Token, state *StateCore) {
	// TODO
}

func Smartquotes(state *StateCore, _ *StateBlock, _ *StateInline, _ int, _ int, _ bool) bool {
	if !state.Md.Options.Typography {
		return false
	}

	for blkIdx := len(state.Tokens) - 1; blkIdx >= 0; blkIdx-- {

		if state.Tokens[blkIdx].Type != "inline" ||
			!QUOTE_TEST_RE.MatchString(state.Tokens[blkIdx].Content) {
			continue
		}

		ProcessInline(state.Tokens[blkIdx].Children, state)
	}
	return true
}
