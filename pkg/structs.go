package pkg

type HighlightFn func(string, string, string) string

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

type LinkReference struct {
	Href  string
	Title string
}

type Env struct {
	info       string
	References map[string]LinkReference
}

type Options struct {
	Html       bool
	XhtmlOut   bool
	Breaks     bool
	Linkify    bool
	LangPrefix string

	// Enable some language-neutral replacements + quotes beautification
	Typography bool
	MaxNesting int
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
