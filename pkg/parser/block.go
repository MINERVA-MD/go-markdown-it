package parser

import (
	. "go-markdown-it/pkg"
	. "go-markdown-it/pkg/rules"
	. "go-markdown-it/pkg/rules/block"
	. "go-markdown-it/pkg/types"
)

type ParserBlock struct {
	Ruler Ruler
}

func (p *ParserBlock) Init() {

}

func (p *ParserBlock) Parse(str string, md *Parser, env Env, outTokens *[]*Token) {

}

func (p *ParserBlock) Tokenize(state *StateBlock, stateLine int, endLine int, silent bool) {
	// TODO
}
