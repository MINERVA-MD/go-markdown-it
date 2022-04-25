package presets

import . "go-markdown-it/pkg/types"

var ZeroPresets = Preset{
	Options: Options{
		Html:       false,
		XhtmlOut:   false,
		Breaks:     false,
		Linkify:    false,
		LangPrefix: "language-",
		Typography: false,
		MaxNesting: 20,
		Quotes: [4]string{
			"\u201c",
			"\u201d",
			"\u2018",
			"\u2019",
		},
		Highlight: nil,
	},
	Components: Components{
		Core: Core{
			Rules: []string{
				"normalize",
				"block",
				"inline",
				"text_join",
			},
		},
		Block: Block{
			Rules: []string{
				"paragraph",
			},
		},
		Inline: Inline{
			Rules: []string{
				"text",
			},
			Rules2: []string{
				"balance_pairs",
				"fragments_join",
			},
		},
	},
}
