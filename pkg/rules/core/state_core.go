package core

import (
	"go-markdown-it/pkg"
	. "go-markdown-it/pkg/types"
)

type StateCore struct {
	Src        string
	Env        Env
	Tokens     []*pkg.Token
	InlineMode bool
	Md         pkg.MarkdownIt
}

func (sc *StateCore) StateCore(src string, md *pkg.MarkdownIt, env Env) {
	sc.Src = src
	sc.Env = env
	sc.Tokens = []*pkg.Token{}
	sc.InlineMode = false
	sc.Md = *md
}
