package pkg

import (
	. "go-markdown-it/pkg/rules/inline"
	. "go-markdown-it/pkg/types"
)

type ParserInline struct{}

type Parser struct {
	Options Options
	Helpers Helpers
	Inline  ParserInline
}

func (p *Parser) NormalizeLink(url string) string {
	// TODO
	return url
}

func (p *Parser) ValidateLink(url string) bool {
	// TODO
	return false
}

func (pi *ParserInline) Tokenize(state *StateInline) {
	// TODO
}
