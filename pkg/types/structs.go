package types

import (
	. "go-markdown-it/pkg/rules"
	. "go-markdown-it/pkg/rules/block"
	. "go-markdown-it/pkg/rules/core"
	. "go-markdown-it/pkg/rules/inline"
	"regexp"
)

type PluginMetaData struct {
	Delimiters []string
}

type Attribute struct {
	Name  string
	Value string
}

type RuleFunction func(*StateCore, *StateBlock, *StateInline, int, int, bool) bool
type HighlightFn func(string, string, string) string
type Cache map[string][]RuleFunction
type Core struct {
	Rules []string
}

type Block struct {
	Rules []string
	Ruler Ruler
}

type Inline struct {
	Rules  []string
	Rules2 []string
}

type Rule struct {
	Name    string
	Enabled bool
	Fn      RuleFunction
	Alt     []string
	// TODO: convert ^this tuple to its own type
}

type Reference struct {
	Href  string
	Title string
}

type Env struct {
	info       string
	References map[string]Reference
}

type Options struct {
	Html       bool
	XhtmlOut   bool
	Breaks     bool
	Linkify    bool
	LangPrefix string

	// Enable some language-neutral replacements + quotes beautification
	Typography bool
	MaxNesting uint8
	Quotes     [4]string
	Highlight  HighlightFn
}

type Components struct {
	Core   Core
	Block  Block
	Inline Inline
}

type Preset struct {
	Options    Options
	Components Components
}

type HtmlSequence struct {
	Start     *regexp.Regexp
	End       *regexp.Regexp
	Terminate bool
}
