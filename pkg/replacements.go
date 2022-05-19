package pkg

import (
	"strings"
)

func ReplaceFn(name string) string {
	return SCOPED_ABBR[strings.ToLower(name)]
}

func ReplaceScoped(inlineTokens []*Token) {
	inside_autolink := 0

	for i := len(inlineTokens) - 1; i >= 0; i-- {
		token := inlineTokens[i]

		if token.Type == "text" && inside_autolink == 0 {
			token.Content = ReplaceAllStringSubmatchFunc(SCOPED_ABBR_RE, token.Content, func(i []string) string {
				return ReplaceFn(i[1])
			})
		}

		if token.Type == "link_open" && token.Info == "auto" {
			inside_autolink--
		}

		if token.Type == "link_close" && token.Info == "auto" {
			inside_autolink++
		}
	}
}

func ReplaceRare(inlineTokens []*Token) {
	insideAutolink := 0

	for i := len(inlineTokens) - 1; i >= 0; i-- {
		token := inlineTokens[i]

		if token.Type == "text" && insideAutolink == 0 {
			if RARE_RE.MatchString(token.Content) {
				token.Content = PLUS_MINUS_RE.ReplaceAllString(token.Content, "±")

				// .., ..., ....... -> …
				// but ?..... & !..... -> ?.. & !..
				token.Content = DOTS3_RE.ReplaceAllString(token.Content, "…")
				token.Content = QE_DOTS3_RE.ReplaceAllString(token.Content, "$1..")

				token.Content = QE4_RE.ReplaceAllString(token.Content, "$1$1$1")
				token.Content = COMMA_RE.ReplaceAllString(token.Content, ",")

				// em-dash
				token.Content, _ = EM_DASH_RE.Replace(token.Content, "$1\u2014", 0, -1)

				// en-dash
				token.Content, _ = EN_DASH1_RE.Replace(token.Content, "$1\u2013", 0, -1)
				token.Content, _ = EN_DASH2_RE.Replace(token.Content, "$1\u2013", 0, -1)
			}
		}

		if token.Type == "link_open" && token.Info == "auto" {
			insideAutolink--
		}

		if token.Type == "link_close" && token.Info == "auto" {
			insideAutolink++
		}
	}
}

func Replace(state *StateCore, _ *StateBlock, _ *StateInline, _ int, _ int, _ bool) bool {

	if !state.Md.Options.Typography {
		return false
	}

	for blkIdx := len(*state.Tokens) - 1; blkIdx >= 0; blkIdx-- {

		if (*state.Tokens)[blkIdx].Type != "inline" {
			continue
		}

		if SCOPED_ABBR_TEST_RE.MatchString((*state.Tokens)[blkIdx].Content) {
			ReplaceScoped((*state.Tokens)[blkIdx].Children)
		}

		if RARE_RE.MatchString((*state.Tokens)[blkIdx].Content) {
			ReplaceRare((*state.Tokens)[blkIdx].Children)
		}
	}
	return true
}
