package pkg

var CommonmarkPresets = Preset{
	Options: Options{
		Html:       true,
		XhtmlOut:   true,
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
				"blockquote",
				"code",
				"fence",
				"heading",
				"hr",
				"html_block",
				"lheading",
				"list",
				"reference",
				"paragraph",
			},
		},
		Inline: Inline{
			Rules: []string{
				"autolink",
				"backticks",
				"emphasis",
				"entity",
				"escape",
				"html_inline",
				"image",
				"link",
				"newline",
				"text",
			},
			Rules2: []string{
				"balance_pairs",
				"emphasis",
				"fragments_join",
			},
		},
	},
}
