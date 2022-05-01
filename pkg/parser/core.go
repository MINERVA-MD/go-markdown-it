package parser

import (
	. "go-markdown-it/pkg/rules"
	. "go-markdown-it/pkg/rules/core"
	"go-markdown-it/pkg/types"
)

type Core struct {
	Ruler Ruler
}

var rules = map[string]types.RuleFunction{
	"normalize":    Normalize,
	"block":        BlockCore,
	"inline":       InlineCore,
	"replacements": Replace,
	"smartquotes":  Smartquotes,
	"text_join":    TextJoin,
}

func (c *Core) Core() {
	c.Ruler = Ruler{
		Rules: []types.Rule{},
		Cache: nil,
	}

	for k, v := range rules {
		c.Ruler.Push(k, v, types.Rule{
			Name:    "",
			Enabled: false,
			Fn:      nil,
			Alt:     nil,
		})
	}
}

func (c *Core) Process(state *StateCore) {
	rules := c.Ruler.GetRules("")

	for _, rule := range rules {
		rule(state, nil, nil, 0, 0, false)
	}
}
