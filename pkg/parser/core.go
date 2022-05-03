package parser

import (
	. "go-markdown-it/pkg/rules"
	. "go-markdown-it/pkg/rules/core"
	"go-markdown-it/pkg/types"
)

type ParserCore struct {
	Ruler Ruler
}

var c_rules = map[string]types.RuleFunction{
	"normalize":    Normalize,
	"block":        BlockCore,
	"inline":       InlineCore,
	"replacements": Replace,
	"smartquotes":  Smartquotes,
	"text_join":    TextJoin,
}

func (c *ParserCore) Core() {
	c.Ruler = Ruler{
		Rules: []types.Rule{},
		Cache: nil,
	}

	for k, v := range c_rules {
		c.Ruler.Push(k, v, types.Rule{
			Name:    k,
			Enabled: false,
			Fn:      v,
			Alt:     nil,
		})
	}
}

func (c *ParserCore) Process(state *StateCore) {
	_rules := c.Ruler.GetRules("")

	for _, rule := range _rules {
		rule(state, nil, nil, 0, 0, false)
	}
}
