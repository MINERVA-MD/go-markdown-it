package pkg

import (
	"unicode/utf8"
)

var i1Rules = []Rule{
	{
		Name:    "text",
		Enabled: false,
		Fn:      Text,
		Alt:     []string{},
	},
	{
		Name:    "linkify",
		Enabled: false,
		Fn:      LLinkify,
		Alt:     []string{},
	},
	{
		Name:    "newline",
		Enabled: false,
		Fn:      Newline,
		Alt:     []string{},
	},
	{
		Name:    "escape",
		Enabled: false,
		Fn:      Escape,
		Alt:     []string{},
	},
	{
		Name:    "backticks",
		Enabled: false,
		Fn:      Backtick,
		Alt:     []string{},
	},
	{
		Name:    "strikethrough",
		Enabled: false,
		Fn:      Strikethrough,
		Alt:     []string{},
	},
	{
		Name:    "emphasis",
		Enabled: false,
		Fn:      Emphasis,
		Alt:     []string{},
	},
	{
		Name:    "link",
		Enabled: false,
		Fn:      Link,
		Alt:     []string{},
	},
	{
		Name:    "image",
		Enabled: false,
		Fn:      Image,
		Alt:     []string{},
	},
	{
		Name:    "autolink",
		Enabled: false,
		Fn:      AutoLink,
		Alt:     []string{},
	},
	{
		Name:    "html_inline",
		Enabled: false,
		Fn:      HtmlInline,
		Alt:     []string{},
	},
	{
		Name:    "entity",
		Enabled: false,
		Fn:      Entity,
		Alt:     []string{},
	},
}

var i2Rules = []Rule{
	{
		Name:    "balance_pairs",
		Enabled: false,
		Fn:      BalancePairs,
		Alt:     []string{},
	},
	{
		Name:    "strikethrough",
		Enabled: false,
		Fn:      SPostProcess,
		Alt:     []string{},
	},
	{
		Name:    "emphasis",
		Enabled: false,
		Fn:      EPostProcess,
		Alt:     []string{},
	},
	// rules for pairs separate '**' into its own text tokens, which may be left unused,
	// rule below merges unused segments back with the rest of the text
	{
		Name:    "fragments_join",
		Enabled: false,
		Fn:      FragmentsJoin,
		Alt:     []string{},
	},
}

type ParserInline struct {
	Ruler  Ruler
	Ruler2 Ruler
}

func (i *ParserInline) ParserInline() {

	i.Ruler = Ruler{
		Rules: []Rule{},
		Cache: nil,
	}

	i.Ruler2 = Ruler{
		Rules: []Rule{},
		Cache: nil,
	}

	for _, rule := range i1Rules {
		i.Ruler.Push(rule.Name, rule.Fn, rule)
	}

	for _, rule := range i2Rules {
		i.Ruler2.Push(rule.Name, rule.Fn, rule)
	}
}

func (i *ParserInline) SkipToken(state *StateInline) {
	var ok bool
	pos := state.Pos
	_rules := i.Ruler.GetRules("")
	maxNesting := state.Md.Options.MaxNesting
	cache := state.Cache

	if _, _ok := cache[pos]; _ok {
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
			ok = rule(nil, nil, state, 0, 0, true)
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
	var idx int
	rules := i.Ruler.GetRules("")
	var n = len(rules)
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
			for idx = 0; idx < n; idx++ {
				ok = rules[idx](nil, nil, state, 0, 0, false)
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

		// TODO: Check this
		cc, _ := state.Src2.CharAt(state.Pos)
		state.Pending += cc
		state.Pos++
	}

	if utf8.RuneCountInString(state.Pending) > 0 {
		state.PushPending()
	}
}

func (i *ParserInline) Parse(src string, md *MarkdownIt, env *Env, outTokens *[]*Token) {
	if utf8.RuneCountInString(src) == 0 {
		return
	}

	state := &StateInline{}
	state.StateInline(src, md, env, outTokens)

	i.Tokenize(state)

	_rules := i.Ruler2.GetRules("")

	for _, rule := range _rules {
		_ = rule(nil, nil, state, 0, 0, false)
	}
}
