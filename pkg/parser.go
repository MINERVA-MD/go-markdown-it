package pkg

import (
	. "go-markdown-it/pkg/parser"
	. "go-markdown-it/pkg/rules/inline"
	. "go-markdown-it/pkg/types"
)

type ParserInline struct{}

type Parser struct {
	Options Options
	Helpers Helpers
	Inline  ParserInline
	Block   ParserBlock
}

func (p *Parser) NormalizeLink(url string) string {
	// TODO
	return url
}

func (p *Parser) ValidateLink(url string) bool {
	// TODO
	return false
}

func (p *Parser) NormalizeLinkText(url string) string {
	// TODO
	return ""
}

func (pi *ParserInline) Tokenize(state *StateInline) {
	// TODO
}

func (pi *ParserInline) Parse(content string, md *Parser, Env Env, tokens *[]*Token) string {
	// TODO

	return ""
}
