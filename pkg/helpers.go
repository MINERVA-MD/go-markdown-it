package pkg

type Helpers struct{}
type LinkResult struct {
	Ok    bool
	Pos   int
	Lines int
	Str   string
}

func (state *Helpers) ParseLinkLabel(start int, disableNested bool) int {
	labelEnd := -1
	// TODO

	return labelEnd
}

func (state *Helpers) ParseLinkDestination(str string, pos int, max int) LinkResult {
	// TODO

	return LinkResult{
		Ok:    false,
		Pos:   0,
		Lines: 0,
		Str:   "",
	}
}

func (state *Helpers) ParseLinkTitle(str string, pos int, max int) LinkResult {
	// TODO

	return LinkResult{
		Ok:    false,
		Pos:   0,
		Lines: 0,
		Str:   "",
	}
}
