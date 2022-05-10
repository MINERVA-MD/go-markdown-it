package pkg

type StateCore struct {
	Src        string
	Env        Env
	Tokens     *[]*Token
	InlineMode bool
	Md         MarkdownIt
}

func (sc *StateCore) StateCore(src string, md *MarkdownIt, env Env) {
	sc.Src = src
	sc.Env = env
	sc.Tokens = &[]*Token{}
	sc.InlineMode = false
	sc.Md = *md
}
