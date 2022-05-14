package pkg

func (state *StateInline) ETokenize(silent bool) bool {
	//fmt.Println("Entered Emphasis Tokenization")
	start := state.Pos
	marker := CharCodeAt(state.Src, start)

	//fmt.Println(start, marker)
	if silent {
		return false
	}

	if marker != 0x5F /* _ */ && marker != 0x2A /* * */ {
		return false
	}

	scanned := state.ScanDelims(state.Pos, marker == 0x2A)

	for i := 0; i < scanned.Length; i++ {
		token := state.Push("text", "", 0)
		token.Content = string(marker)

		//fmt.Println(token.Content)
		*state.Delimiters = append(*state.Delimiters, &Delimiter{
			Marker: marker,
			Length: scanned.Length,
			Token:  len(*state.Tokens) - 1,
			End:    -1,
			Open:   scanned.CanOpen,
			Close:  scanned.CanClose,
		})
	}

	state.Pos += scanned.Length
	return true
}

func (state *StateInline) _PostProcess(delim string, idx int) {
	var i int
	var startDelim *Delimiter
	var delimiters *[]*Delimiter

	if delim == "delimiters" {
		delimiters = state.Delimiters
	} else {
		delimiters = state.TokensMeta[idx].Delimiters
	}

	//fmt.Println("Entered Emphasis Post Process")
	max := len(*delimiters)

	//utils.PrettyPrint(delimiters)
	for i = max - 1; i >= 0; i-- {
		startDelim = (*delimiters)[i]

		if startDelim.Marker != 0x5F /* _ */ && startDelim.Marker != 0x2A /* * */ {
			continue
		}

		// Process only opening markers
		if startDelim.End == -1 {
			continue
		}

		endDelim := (*delimiters)[startDelim.End]

		// If the previous delimiter has the same marker and is adjacent to this one,
		// merge those into one strong delimiter.
		//
		// `<em><em>whatever</em></em>` -> `<strong>whatever</strong>`
		//utils.PrettyPrint(startDelim)
		isStrong := i > 0 &&
			(*delimiters)[i-1].End == startDelim.End+1 &&
			// check that first two markers match and adjacent
			(*delimiters)[i-1].Marker == startDelim.Marker &&
			(*delimiters)[i-1].Token == startDelim.Token-1 &&
			// check that last two markers are adjacent (we can safely assume they match)
			(*delimiters)[startDelim.End+1].Token == endDelim.Token+1

		ch := string(startDelim.Marker)

		token := (*state.Tokens)[startDelim.Token]

		if isStrong {
			token.Type = "strong_open"
			token.Tag = "strong"
			token.Markup = ch + ch
		} else {
			token.Type = "em_open"
			token.Tag = "em"
			token.Markup = ch
		}

		token.Nesting = 1
		token.Content = ""

		token = (*state.Tokens)[endDelim.Token]
		if isStrong {
			token.Type = "strong_close"
			token.Tag = "strong"
			token.Markup = ch + ch
		} else {
			token.Type = "em_close"
			token.Tag = "em"
			token.Markup = ch
		}

		token.Nesting = -1
		token.Content = ""

		if isStrong {
			(*state.Tokens)[(*delimiters)[i-1].Token].Content = ""
			(*state.Tokens)[(*delimiters)[startDelim.End+1].Token].Content = ""
			i--
		}
	}
}

func Emphasis(
	_ *StateCore,
	_ *StateBlock,
	state *StateInline,
	_ int,
	_ int,
	silent bool,
) bool {
	return state.ETokenize(silent)
}

func EPostProcess(
	_ *StateCore,
	_ *StateBlock,
	state *StateInline,
	_ int,
	_ int,
	_ bool,
) bool {
	//fmt.Println("Processing Emphasis")
	state.EPostProcess()
	return true
}

func (state *StateInline) EPostProcess() {
	tokensMeta := state.TokensMeta

	//utils.PrettyPrint(state.Delimiters)
	state._PostProcess("delimiters", -1)

	for idx, tokenMeta := range tokensMeta {
		if tokenMeta.Delimiters != nil && len(*tokenMeta.Delimiters) > 0 {
			state._PostProcess("metaDelimiters", idx)
		}
	}
}
