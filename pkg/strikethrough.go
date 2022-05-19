package pkg

func (state *StateInline) Tokenize(silent bool) bool {

	start := state.Pos
	marker, _ := state.Src2.CharCodeAt(start)

	if silent {
		return false
	}

	if marker != 0x7E {
		return false
	}

	scanned := state.ScanDelims(state.Pos, true)
	n := scanned.Length
	ch := string(marker)

	if n < 2 {
		return false
	}

	if n%2 > 0 {
		token := state.Push("text", "", 0)
		token.Content = ch + ch

		*state.Delimiters = append(*state.Delimiters, &Delimiter{
			Marker: marker,
			Length: 0,
			Token:  len(*state.Tokens) - 1,
			End:    -1,
			Open:   scanned.CanOpen,
			Close:  scanned.CanClose,
		})
	}

	state.Pos += scanned.Length

	return true
}

func (state *StateInline) PostProcess(delim string, idx int) {

	var loneMarkers []int
	var delimiters *[]*Delimiter

	if delim == "delimiters" {
		delimiters = state.Delimiters
	} else {
		delimiters = state.TokensMeta[idx].Delimiters
	}

	for _, delimiter := range *delimiters {

		if delimiter.Marker != 0x7E {
			continue
		}

		if delimiter.End == -1 {
			continue
		}

		endDelim := (*delimiters)[delimiter.End]

		token := (*state.Tokens)[delimiter.Token]
		token.Type = "s_open"
		token.Tag = "s"
		token.Nesting = 1
		token.Markup = "~~"
		token.Content = ""

		token = (*state.Tokens)[endDelim.Token]
		token.Type = "s_close"
		token.Tag = "s"
		token.Nesting = -1
		token.Markup = "~~"
		token.Content = ""

		if (*state.Tokens)[endDelim.Token-1].Type == "text" &&
			(*state.Tokens)[endDelim.Token-1].Content == "~" {
			loneMarkers = append(loneMarkers, endDelim.Token-1)
		}
	}

	// If a marker sequence has an odd number of characters, it's splitted
	// like this: `~~~~~` -> `~` + `~~` + `~~`, leaving one marker at the
	// start of the sequence.
	//
	// So, we have to move all those markers after subsequent s_close tags.

	var i int
	var j int
	for len(loneMarkers) > 0 {
		i, loneMarkers = Pop(loneMarkers)
		j = i + 1

		for j < len(*state.Tokens) &&
			(*state.Tokens)[j].Type == "s_close" {
			j++
		}

		j--

		if i != j {
			token := (*state.Tokens)[j]
			(*state.Tokens)[j] = (*state.Tokens)[i]
			(*state.Tokens)[i] = token
		}
	}
}

func Strikethrough(
	_ *StateCore,
	_ *StateBlock,
	state *StateInline,
	_ int,
	_ int,
	silent bool,
) bool {
	return state.Tokenize(silent)
}

func SPostProcess(
	_ *StateCore,
	_ *StateBlock,
	state *StateInline,
	_ int,
	_ int,
	_ bool,
) bool {
	state.Strikethrough()
	return true
}

func (state *StateInline) Strikethrough() {
	tokensMeta := state.TokensMeta
	state.PostProcess("delimiters", -1)

	for idx, tokenMeta := range tokensMeta {
		if tokenMeta.Delimiters != nil && len(*tokenMeta.Delimiters) > 0 {
			state.PostProcess("metaDelimiters", idx)
			//(tokenMeta.Delimiters)
		}
	}
}
