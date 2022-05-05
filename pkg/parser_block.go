package pkg

var b_rules = []Rule{
	{
		Name:    "table",
		Enabled: false,
		Fn:      Table,
		Alt:     []string{"paragraph", "reference"},
	},
	{
		Name:    "code",
		Enabled: false,
		Fn:      Code,
		Alt:     []string{},
	},
	{
		Name:    "fence",
		Enabled: false,
		Fn:      Fence,
		Alt:     []string{"paragraph", "reference", "blockquote", "list"},
	},
	{
		Name:    "blockquote",
		Enabled: false,
		Fn:      BlockQuote,
		Alt:     []string{"paragraph", "reference", "blockquote", "list"},
	},
	{
		Name:    "hr",
		Enabled: false,
		Fn:      Hr,
		Alt:     []string{"paragraph", "reference", "blockquote", "list"},
	},
	{
		Name:    "list",
		Enabled: false,
		Fn:      List,
		Alt:     []string{"paragraph", "reference", "blockquote"},
	},
	{
		Name:    "reference",
		Enabled: false,
		Fn:      Reference,
		Alt:     []string{},
	},
	{
		Name:    "html_block",
		Enabled: false,
		Fn:      HtmlBlock,
		Alt:     []string{"paragraph", "reference", "blockquote"},
	},
	{
		Name:    "heading",
		Enabled: false,
		Fn:      Heading,
		Alt:     []string{"paragraph", "reference", "blockquote"},
	},
	{
		Name:    "lheading",
		Enabled: false,
		Fn:      LHeading,
		Alt:     []string{},
	},
	{
		Name:    "paragraph",
		Enabled: false,
		Fn:      Paragraph,
		Alt:     []string{},
	},
}

type ParserBlock struct {
	Ruler Ruler
}

func (b *ParserBlock) ParserBlock() {
	b.Ruler = Ruler{
		Rules: []Rule{},
		Cache: nil,
	}

	for _, rule := range b_rules {
		b.Ruler.Push(rule.Name, rule.Fn, rule)
	}
}

func (b *ParserBlock) Tokenize(state *StateBlock, startLine int, endLine int, silent bool) {
	var ok bool
	rules := b.Ruler.GetRules("")
	_len := len(rules)
	line := startLine
	hasEmptyLines := false
	maxNesting := state.Md.Options.MaxNesting

	for line < endLine {
		state.Line = state.SkipEmptyLines(line)
		line = state.Line

		if line >= endLine {
			break
		}

		// Termination condition for nested calls.
		// Nested calls currently used for blockquotes & lists
		if state.SCount[line] < state.BlkIndent {
			break
		}

		// If nesting level exceeded - skip tail to the end. That's not ordinary
		// situation and we should not care about content.
		if state.Level >= maxNesting {
			state.Line = endLine
			break
		}

		// Try all possible rules.
		// On success, rule should:
		//
		// - update `state.line`
		// - update `state.tokens`
		// - return true
		for i := 0; i < _len; i++ {
			ok = rules[i](nil, state, nil, line, endLine, false)
			if ok {
				break
			}
		}

		// set state.tight if we had an empty line before current tag
		// i.e. latest empty line should not count
		state.Tight = !hasEmptyLines

		// paragraph might "eat" one newline after it in nested lists
		if state.IsEmpty(state.Line - 1) {
			hasEmptyLines = true
		}

		line = state.Line

		if line < endLine && state.IsEmpty(line) {
			hasEmptyLines = true
			line++
			state.Line = line
		}
	}
}

func (b *ParserBlock) Parse(src string, md *MarkdownIt, env Env, outTokens *[]*Token) {
	if len(src) == 0 {
		return
	}

	state := &StateBlock{}
	state.StateBlock(src, md, env, outTokens)

	b.Tokenize(state, state.Line, state.LineMax, false)
}
