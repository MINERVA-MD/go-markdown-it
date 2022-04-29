package parser

import (
	. "go-markdown-it/pkg/rules"
	. "go-markdown-it/pkg/rules/block"
)

type ParserBlock struct {
	Ruler Ruler
}

func (p *ParserBlock) Init() {

}

func (p *ParserBlock) Tokenize(state *StateBlock, stateLine int, endLine int, silent bool) {
	// TODO
}
