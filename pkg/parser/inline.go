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

type ParserInline struct {
	Ruler  rules.Ruler
	Ruler2 rules.Ruler
}

func (i *ParserInline) ParserInline() {

	i.Ruler = rules.Ruler{
		Rules: []types.Rule{},
		Cache: nil,
	}

	i.Ruler2 = rules.Ruler{
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
		i.Ruler2.Push(k, v, types.Rule{
			Name:    k,
			Enabled: false,
			Fn:      v,
			Alt:     nil,
		})
	}
}

func (i *ParserInline) SkipToken(state *StateInline) {
	var ok bool
	pos := state.Pos
	_rules := i.Ruler.GetRules("")
	maxNesting := state.Md.Options.MaxNesting
	cache := state.Cache

	if _, ok := cache[pos]; ok {
		state.Pos = cache[pos]
		return
	}

	if state.Level < maxNesting {
		for _, rule := range _rules {
			// Increment state.level and decrement it later to limit recursion.
			// It's harmless to do here, because no tokens are created. But ideally,
			// we'd need a separate private state variable for this purpose.
			//
			state.Level++
			ok = rule(nil, nil, state, 0, 0, false)
			state.Level--

			if ok {
				break
			}
		}
	} else {
		// Too much nesting, just skip until the end of the paragraph.
		//
		// NOTE: this will cause links to behave incorrectly in the following case,
		//       when an amount of `[` is exactly equal to `maxNesting + 1`:
		//
		//       [[[[[[[[[[[[[[[[[[[[[foo]()
		//
		// TODO: remove this workaround when CM standard will allow nested links
		//       (we can replace it by preventing links from being parsed in
		//       validation mode)
		//
		state.Pos = state.PosMax
	}

	if !ok {
		state.Pos++
	}
	cache[pos] = state.Pos
}

func (i *ParserInline) Tokenize(state *StateInline) {
	_rules := i.Ruler.GetRules("")
	end := state.PosMax
	maxNesting := state.Md.Options.MaxNesting

	for state.Pos < end {
		// Try all possible rules.
		// On success, rule should:
		//
		// - update `state.pos`
		// - update `state.tokens`
		// - return true

		var ok bool
		if state.Level < maxNesting {
			for _, rule := range _rules {
				ok = rule(nil, nil, state, 0, 0, false)
				if ok {
					break
				}
			}
		}

		if ok {
			if state.Pos >= end {
				break
			}
			continue
		}

		state.Pending += string(state.Src[state.Pos])
		state.Pos++
	}
	if len(state.Pending) > 0 {
		state.PushPending()
	}
}

func (i *ParserInline) Parse(src string, md *pkg.MarkdownIt, env types.Env, outTokens []*pkg.Token) {
	if len(src) == 0 {
		return
	}

	state := &StateInline{}
	state.StateInline(src, md, env, outTokens)

	i.Tokenize(state)

	_rules := i.Ruler2.GetRules("")

	for _, rule := range _rules {
		rule(nil, nil, state, 0, 0, false)
	}
}
