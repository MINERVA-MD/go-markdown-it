package pkg

var DefaultPresets = Preset{
	Options: Options{
		Html:       false,
		XhtmlOut:   false,
		Breaks:     false,
		Linkify:    false,
		LangPrefix: "language-",
		Typography: false,
		MaxNesting: 100,
		Quotes: [4]string{
			"\u201c",
			"\u201d",
			"\u2018",
			"\u2019",
		},
		Highlight: nil,
	},
	Components: Components{
		Core:   Core{},
		Block:  Block{},
		Inline: Inline{},
	},
}
