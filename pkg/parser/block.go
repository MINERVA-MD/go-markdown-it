package parser

import (
	. "go-markdown-it/pkg"
	. "go-markdown-it/pkg/rules"
	. "go-markdown-it/pkg/rules/block"
	"go-markdown-it/pkg/types"
)

var b_rules = map[string]types.Rule{
	"table": {
		Name:    "table",
		Enabled: false,
		Fn:      Table,
		Alt:     []string{"paragraph", "reference"},
	},
	"code": {
		Name:    "code",
		Enabled: false,
		Fn:      Code,
		Alt:     []string{},
	},
	"fence": {
		Name:    "fence",
		Enabled: false,
		Fn:      Fence,
		Alt:     []string{"paragraph", "reference", "blockquote", "list"},
	},
	"blockquote": {
		Name:    "blockquote",
		Enabled: false,
		Fn:      BlockQuote,
		Alt:     []string{"paragraph", "reference", "blockquote", "list"},
	},
	"hr": {
		Name:    "hr",
		Enabled: false,
		Fn:      Hr,
		Alt:     []string{"paragraph", "reference", "blockquote", "list"},
	},
	"list": {
		Name:    "list",
		Enabled: false,
		Fn:      List,
		Alt:     []string{"paragraph", "reference", "blockquote"},
	},
	"reference": {
		Name:    "reference",
		Enabled: false,
		Fn:      Reference,
		Alt:     []string{},
	},
	"html_block": {
		Name:    "html_block",
		Enabled: false,
		Fn:      HtmlBlock,
		Alt:     []string{"paragraph", "reference", "blockquote"},
	},
	"heading": {
		Name:    "heading",
		Enabled: false,
		Fn:      Heading,
		Alt:     []string{"paragraph", "reference", "blockquote"},
	},
	"lheading": {
		Name:    "lheading",
		Enabled: false,
		Fn:      LHeading,
		Alt:     []string{},
	},
	"paragraph": {
		Name:    "paragraph",
		Enabled: false,
		Fn:      Paragraph,
		Alt:     []string{},
	},
}

type Block struct {
	Ruler Ruler
}

func (b *Block) ParserBlock() {
	b.Ruler = Ruler{
		Rules: []types.Rule{},
		Cache: nil,
	}

	for k, rule := range b_rules {
		b.Ruler.Push(k, rule.Fn, rule)
	}
}

func (b *Block) Tokenize(state *StateBlock, startLine int, endLine int, silent bool) {
	var ok bool
	rules := b.Ruler.GetRules("")
	_len := len(rules)
	line := startLine
	hasEmptyLines := false
	maxNesting := state.Md.Options.MaxNesting

	for line < endLine {
		line = state.SkipEmptyLines(line)
		state.Line = state.SkipEmptyLines(line)
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

func (b *Block) Parse(str string, md *Parser, env types.Env, outTokens []*Token) {
	if len(str) == 0 {
		return
	}

	// TODO
	state := StateBlockInit()
	b.Tokenize(state, state.Line, state.LineMax, false)
}
