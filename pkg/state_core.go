package pkg

type StateCore struct {
	Src        string
	Src2       *MDString
	Env        *Env
	Tokens     *[]*Token
	InlineMode bool
	Md         *MarkdownIt
}

func (sc *StateCore) StateCore(src string, md *MarkdownIt, env *Env) {
	mds := &MDString{}
	_ = mds.Init(src)

	sc.Src = src
	sc.Src2 = mds
	sc.Env = env
	sc.Tokens = &[]*Token{}
	sc.InlineMode = false
	sc.Md = md
}
