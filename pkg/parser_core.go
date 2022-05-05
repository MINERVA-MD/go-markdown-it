package pkg

type ParserCore struct {
	Ruler Ruler
}

var c_rules = map[string]RuleFunction{
	"normalize":    Normalize,
	"block":        BlockCore,
	"inline":       InlineCore,
	"replacements": Replace,
	"smartquotes":  Smartquotes,
	"text_join":    TextJoin,
}

func (c *ParserCore) ParserCore() {
	c.Ruler = Ruler{
		Rules: []Rule{},
		Cache: nil,
	}

	for k, v := range c_rules {
		c.Ruler.Push(k, v, Rule{
			Name:    k,
			Enabled: false,
			Fn:      v,
			Alt:     nil,
		})
	}
}

func (c *ParserCore) Process(state *StateCore) {
	_rules := c.Ruler.GetRules("")

	for idx := 0; idx < len(_rules); idx++ {
		_rules[idx](state, nil, nil, 0, 0, false)
	}
}
