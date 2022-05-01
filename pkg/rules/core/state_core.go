package core

import (
	. "go-markdown-it/pkg"
	. "go-markdown-it/pkg/types"
)

type StateCore struct {
	Src        string
	Env        Env
	Tokens     []*Token
	InlineMode bool
	Md         Parser
}

func (sc *StateCore) StateCore(src string, md *Parser, env Env) {
	sc.Src = src
	sc.Env = env
	sc.Tokens = []*Token{}
	sc.InlineMode = false
	sc.Md = *md
}
