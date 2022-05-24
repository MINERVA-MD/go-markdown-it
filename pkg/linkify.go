package pkg

import (
	"gitlab.com/golang-commonmark/linkify"
	"unicode/utf8"
)

type Linkify struct{}

//func ILinkify(
//	_ *StateCore,
//	_ *StateBlock,
//	state *StateInline,
//	_ int,
//	_ int,
//	silent bool,
//) bool {
//	return state.ILinkify(silent)
//}

func (l *Linkify) Test(_ string) bool {
	return false
}

func (l *Linkify) Pretest(_ string) bool {
	return false
}

func LLinkify(
	_ *StateCore,
	_ *StateBlock,
	state *StateInline,
	_ int,
	_ int,
	silent bool,
) bool {
	return state.Linkify(silent)
}

func (state *StateInline) Linkify(silent bool) bool {

	if !state.Md.Options.Linkify {
		return false
	}

	if state.LinkLevel > 0 {
		return false
	}

	pos := state.Pos
	max := state.PosMax

	if pos+3 > max {
		return false
	}

	if cc1, _ := state.Src2.CharCodeAt(pos); cc1 != 0x3A {
		return false
	}
	if cc2, _ := state.Src2.CharCodeAt(pos + 1); cc2 != 0x2F {
		return false
	}
	if cc3, _ := state.Src2.CharCodeAt(pos + 2); cc3 != 0x2F {
		return false
	}

	match := SCHEME_RE.FindStringSubmatch(state.Pending2.String())
	proto := match[1]

	// 	link := state.Md.LLinkify.matchAtStart(state.src.slice(pos - proto.length));
	// TODO: Make proper call ^
	link := state.Src2.Slice(pos-utf8.RuneCountInString(proto), state.Src2.Length)
	if utf8.RuneCountInString(link) == 0 {
		return false
	}

	// TODO: url = link.url
	url := link
	url = LINKIFY_CONFLICT_RE.ReplaceAllString(url, "")

	fullUrl := state.Md.NormalizeLink(url)
	if !state.Md.ValidateLink(fullUrl) {
		return false
	}

	if !silent {
		// TODO: double check negative start!!!!!
		l := utf8.RuneCountInString(proto)
		_ = state.Pending2.Init(state.Pending2.Slice(0, state.Pending2.Length-l))

		token := state.Push("link_open", "a", 1)
		token.Attrs = []Attribute{
			{
				Name:  "href",
				Value: fullUrl,
			},
		}

		token.Markup = "linkify"
		token.Info = "auto"

		token = state.Push("text", "", 0)
		token.Content = state.Md.NormalizeLinkText(url)

		token = state.Push("link_close", "a", -1)
		token.Markup = "linkify"
		token.Info = "auto"
	}

	state.Pos += utf8.RuneCountInString(url) - utf8.RuneCountInString(proto)

	return true
}

func ILinkify(
	_ *StateCore,
	_ *StateBlock,
	state *StateInline,
	_ int,
	_ int,
	_ bool,
) bool {

	var pos int
	var level int
	var url string
	var lastPos int
	var text string
	var urlText string
	var fullUrl string
	var token Token
	var htmlLinkLevel int
	var nodes []*Token
	var tokens []*Token
	var links []linkify.Link
	var currentToken *Token

	blockTokens := state.Tokens

	if !state.Md.Options.Linkify {
		return false
	}

	l := len(*blockTokens)
	for j := 0; j < l; j++ {
		if (*blockTokens)[j].Type != "inline" ||
			!state.Md.Linkify.Pretest((*blockTokens)[j].Content) {
			continue
		}

		tokens = (*blockTokens)[j].Children

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

				nodes = []*Token{}
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

					urlText = Slice(text, link.Start, link.End)

					// Linkifier might send raw hostnames like "example.com", where url
					// starts with domain name. So we prepend http:// in those cases,
					// and remove it afterwards.
					//
					if utf8.RuneCountInString(link.Scheme) == 0 {
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
						token = GenerateToken("text", "", 0)
						token.Content = Slice(text, lastPos, pos)
						token.Level = level
						nodes = append(nodes, &token)
					}

					token = GenerateToken("link_open", "a", 1)
					token.Attrs = []Attribute{
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

					token = GenerateToken("text", "", 0)

					token.Content = urlText
					token.Level = level
					nodes = append(nodes, &token)

					token = GenerateToken("link_close", "a", -1)

					level--
					token.Level = level
					token.Markup = "linkify"
					token.Info = "auto"
					nodes = append(nodes, &token)

					lastPos = link.End
				}
				if lastPos < utf8.RuneCountInString(text) {
					token = GenerateToken("text", "", 0)
					token.Content = Slice(text, lastPos, utf8.RuneCountInString(text))
					//text[lastPos:]
					token.Level = level
					nodes = append(nodes, &token)
				}

				tokens = InsertTokensAt(tokens, i, nodes)
				(*blockTokens)[j].Children = tokens
			}
		}
	}
	return true
}

func InsertTokensAt(a []*Token, index int, values []*Token) []*Token {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, values...)
	}
	one := append(a[:index], values...)
	a = append(one, a[index+1:]...)
	return a
}
