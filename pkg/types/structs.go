package types

type PluginMetaData struct {
	Value string
}

type Attribute struct {
	Name  string
	Value string
}

type Token struct {
	// Type of the token (string, e.g. "paragraph_open")
	Type string

	// HTML Tag name, e.g. "p"
	Tag string

	// Html attributes. Format: `[ [ name1, value1 ], [ name2, value2 ]
	Attrs []Attribute

	// Source map info. Format: `[ line_begin, line_end ]`
	Map []string

	/**
	 * Level change (number in {-1, 0, 1} set), where:
	 *
	 * -  `1` means the tag is opening
	 * -  `0` means the tag is self-closing
	 * - `-1` means the tag is closing
	 **/
	Nesting int8

	// nesting level, the same as `state.level
	Level int8

	// An array of child nodes (inline and img tokens)
	Children []*Token

	// In a case of self-closing tag (code, html, fence, etc.),
	// it has contents of this tag.
	Content string

	// '*' or '_' for emphasis, fence string for fence, etc.
	Markup string

	/**
	 * Additional information:
	 *
	 * - Info string for "fence" tokens
	 * - The value "auto" for autolink "link_open" and "link_close" tokens
	 * - The string value of the item marker for ordered-list "list_item_open" tokens
	 **/
	Info string

	// A place for plugins to store an arbitrary data
	Meta PluginMetaData

	// True for block-level tokens, false for inline tokens.
	//  Used in renderer to calculate line breaks
	Block bool

	// If it's true, ignore this element when rendering. Used for tight lists
	// to hide paragraphs.
	Hidden bool
}

type RuleFunction func()
type HighlightFn func(string, string, string) string
type Cache map[string][]RuleFunction
type Core struct {
	Rules []string
}

type Block struct {
	Rules []string
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

type Env struct {
	info string
}

type Ruler struct {
	Rules []Rule
	Cache Cache
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

//type Renderer interface {
//	RenderAttrs(token *Token) string
//	RenderToken(tokens []*Token, idx int, options Options)
//	RenderInline(tokens []*Token, options Options, env string)
//	RenderInlineAsText()
//}
