package core

import (
	. "go-markdown-it/pkg"
	. "go-markdown-it/pkg/common"
	. "go-markdown-it/pkg/maps"
	. "go-markdown-it/pkg/types"
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
	inside_autolink := 0

	for i := len(inlineTokens) - 1; i >= 0; i-- {
		token := inlineTokens[i]

		if token.Type == "text" && inside_autolink == 0 {
			if RARE_RE.MatchString(token.Content) {
				token.Content = PLUS_MINUS_RE.ReplaceAllString(token.Content, "±")

				// .., ..., ....... -> …
				// but ?..... & !..... -> ?.. & !..
				token.Content = DOTS3_RE.ReplaceAllString(token.Content, "…")
				token.Content = QE_DOTS3_RE.ReplaceAllString(token.Content, "$1..")

				token.Content = QE4_RE.ReplaceAllString(token.Content, "$1$1$1")
				token.Content = COMMA_RE.ReplaceAllString(token.Content, ",")

				// em-dash
				token.Content = EM_DASH_RE.ReplaceAllString(token.Content, "$1\\u2014")

				// en-dash
				token.Content = EN_DASH1_RE.ReplaceAllString(token.Content, "$1\\u2013")
				token.Content = EN_DASH2_RE.ReplaceAllString(token.Content, "$1\\u2013")

			}
		}

		if token.Type == "link_open" && token.Info == "auto" {
			inside_autolink--
		}

		if token.Type == "link_close" && token.Info == "auto" {
			inside_autolink++
		}
	}
}

func Replace(state *StateCore) {

	if !state.Md.Options.Typography {
		return
	}

	for blkIdx := len(state.Tokens) - 1; blkIdx >= 0; blkIdx-- {

		if state.Tokens[blkIdx].Type != "inline" {
			continue
		}

		if SCOPED_ABBR_TEST_RE.MatchString(state.Tokens[blkIdx].Content) {
			ReplaceScoped(state.Tokens[blkIdx].Children)
		}

		if RARE_RE.MatchString(state.Tokens[blkIdx].Content) {
			ReplaceRare(state.Tokens[blkIdx].Children)
		}

	}
}
