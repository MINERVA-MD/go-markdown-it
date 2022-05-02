package parser

import (
	"go-markdown-it/pkg"
	"go-markdown-it/pkg/rules"
	. "go-markdown-it/pkg/rules/inline"
	"go-markdown-it/pkg/types"
)

var i1_rules = map[string]types.RuleFunction{
	"text":          Text,
	"linkify":       Linkify,
	"newline":       Newline,
	"escape":        Escape,
	"backticks":     Backtick,
	"strikethrough": Strikethrough,
	"emphasis":      Emphasis,
	"link":          Link,
	"image":         Image,
	"autolink":      AutoLink,
	"html_inline":   HtmlInline,
	"entity":        Entity,
}

var i2_rules = map[string]types.RuleFunction{
	"balance_pairs": BalancePairs,
	"strikethrough": SPostProcess,
	"emphasis":      EPostProcess,
	// rules for pairs separate '**' into its own text tokens, which may be left unused,
	// rule below merges unused segments back with the rest of the text
	"fragments_join": FragmentsJoin,
}

type Inline struct {
	Ruler rules.Ruler
}

func (i *Inline) ParserInline() {

	i.Ruler = rules.Ruler{
		Rules: []types.Rule{},
		Cache: nil,
	}

	for k, v := range i1_rules {
		i.Ruler.Push(k, v, types.Rule{
			Name:    k,
			Enabled: false,
			Fn:      v,
			Alt:     nil,
		})
	}

	for k, v := range i2_rules {
		i.Ruler.Push(k, v, types.Rule{
			Name:    k,
			Enabled: false,
			Fn:      v,
			Alt:     nil,
		})
	}
}

func (i *Inline) SkipToken() {
	// TODO
}

func (i *Inline) Tokenize() {
	// TODO
}

func (i *Inline) Parse(str string, md *pkg.Parser, env types.Env, outTokens []*pkg.Token) {
	if len(str) == 0 {
		return
	}

	// TODO
}
