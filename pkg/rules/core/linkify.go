package core

import (
	"gitlab.com/golang-commonmark/linkify"
	"go-markdown-it/pkg"
	. "go-markdown-it/pkg/common"
	. "go-markdown-it/pkg/rules/block"
	. "go-markdown-it/pkg/rules/inline"
	"go-markdown-it/pkg/types"
)

func Linkify(state *StateCore, _ *StateBlock, _ *StateInline, _ int, _ int, _ bool) bool {

	var pos int
	var level int
	var url string
	var lastPos int
	var text string
	var urlText string
	var fullUrl string
	var token pkg.Token
	var htmlLinkLevel int
	var nodes []*pkg.Token
	var tokens []*pkg.Token
	var links []linkify.Link
	var currentToken *pkg.Token
	blockTokens := state.Tokens

	if !state.Md.Options.Linkify {
		return false
	}

	l := len(blockTokens)
	for j := 0; j < l; j++ {
		if blockTokens[j].Type != "inline" ||
			!state.Md.Linkify.Pretest(blockTokens[j].Content) {
			continue
		}

		tokens = blockTokens[j].Children

		htmlLinkLevel = 0

		// We scan from the end, to keep position when new tags added.
		// Use reversed logic in links start/end match
		for i := len(tokens) - 1; i >= 0; i-- {
			currentToken = tokens[i]

			// Skip content of markdown links
			if currentToken.Type == "link_close" {
				i--
				for tokens[i].Level != currentToken.Level && tokens[i].Type != "link_open" {
					i--
				}
				continue
			}

			// Skip content of html tag links
			if currentToken.Type == "html_inline" {
				if IsLinkOpen(currentToken.Content) && htmlLinkLevel > 0 {
					htmlLinkLevel--
				}
				if IsLinkClose(currentToken.Content) {
					htmlLinkLevel++
				}
			}

			if htmlLinkLevel > 0 {
				continue
			}

			if currentToken.Type == "text" && state.Md.Linkify.Test(currentToken.Content) {
				text = currentToken.Content
				links = linkify.Links(text)

				nodes = []*pkg.Token{}
				level = currentToken.Level
				lastPos = 0

				// forbid escape sequence at the start of the string,
				// this avoids http\://example.com/ from being linkified as
				// http:<a href="//example.com/">//example.com/</a>
				if len(links) > 0 &&
					links[0].Start == 0 &&
					i > 0 &&
					tokens[i-1].Type == "text_special" {
					links = links[1:]
				}

				for ln := 0; ln < len(links); ln++ {
					link := links[ln]
					url = link.Scheme + text[link.Start:link.End]
					fullUrl = state.Md.NormalizeLink(url)
					if !state.Md.ValidateLink(fullUrl) {
						continue
					}

					urlText = text[link.Start:link.End]

					// Linkifier might send raw hostnames like "example.com", where url
					// starts with domain name. So we prepend http:// in those cases,
					// and remove it afterwards.
					//
					if len(link.Scheme) == 0 {
						urlText = state.Md.NormalizeLinkText("http://" + urlText)
						urlText = HTTP_RE.ReplaceAllString(urlText, "")
					} else if link.Scheme == "mailto:" && !MAILTO_RE.MatchString(urlText) {
						urlText = state.Md.NormalizeLinkText("mailto:" + urlText)
						urlText = MAILTO_RE.ReplaceAllString(urlText, "")
					} else {
						urlText = state.Md.NormalizeLinkText(urlText)
					}

					pos = link.Start

					if pos > lastPos {
						token = pkg.GenerateToken("text", "", 0)
						token.Content = text[lastPos:pos]
						token.Level = level
						nodes = append(nodes, &token)
					}

					token = pkg.GenerateToken("link_open", "a", 1)
					token.Attrs = []types.Attribute{
						{
							Name:  "href",
							Value: fullUrl,
						},
					}

					token.Level = level
					level++

					token.Markup = "linkify"
					token.Info = "auto"
					nodes = append(nodes, &token)

					token = pkg.GenerateToken("text", "", 0)
					token.Content = urlText
					token.Level = level
					nodes = append(nodes, &token)

					token = pkg.GenerateToken("link_close", "a", -1)

					level--
					token.Level = level
					token.Markup = "linkify"
					token.Info = "auto"
					nodes = append(nodes, &token)

					lastPos = link.End
				}
				if lastPos < len(text) {
					token = pkg.GenerateToken("text", "", 0)
					token.Content = text[lastPos:]
					token.Level = level
					nodes = append(nodes, &token)
				}

				tokens = InsertTokensAt(tokens, i, nodes)
				blockTokens[j].Children = tokens
			}
		}
	}

	return true
}

func InsertTokensAt(a []*pkg.Token, index int, values []*pkg.Token) []*pkg.Token {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, values...)
	}
	one := append(a[:index], values...)
	a = append(one, a[index+1:]...)
	return a
}
