package pkg

type ParserCore struct {
	Ruler Ruler
}

var cRules = []Rule{
	{
		Name:    "normalize",
		Enabled: false,
		Fn:      Normalize,
		Alt:     []string{},
	},
	{
		Name:    "block",
		Enabled: false,
		Fn:      BlockCore,
		Alt:     []string{},
	},
	{
		Name:    "inline",
		Enabled: false,
		Fn:      InlineCore,
		Alt:     []string{},
	},
	{
		Name:    "replacements",
		Enabled: false,
		Fn:      Replace,
		Alt:     []string{},
	},
	{
		Name:    "smartquotes",
		Enabled: false,
		Fn:      Smartquotes,
		Alt:     []string{},
	},
	{
		Name:    "text_join",
		Enabled: false,
		Fn:      TextJoin,
		Alt:     []string{},
	},
}

func (c *ParserCore) ParserCore() {
	c.Ruler = Ruler{
		Rules: []Rule{},
		Cache: nil,
	}

	for _, rule := range cRules {
		c.Ruler.Push(rule.Name, rule.Fn, rule)
	}
}

func (c *ParserCore) Process(state *StateCore) {
	rules := c.Ruler.GetRules("")

	for idx := 0; idx < len(rules); idx++ {
		_ = rules[idx](state, nil, nil, 0, 0, false)
	}
}
