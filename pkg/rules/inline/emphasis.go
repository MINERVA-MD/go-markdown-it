package inline

import (
	. "go-markdown-it/pkg/rules/block"
	"go-markdown-it/pkg/rules/core"
)

func (state *StateInline) ETokenize(silent bool) bool {
	start := state.Pos
	marker := CharCodeAt(state.Src, start)

	if silent {
		return false
	}

	if marker != 0x5F /* _ */ && marker != 0x2A /* * */ {
		return false
	}

	scanned := state.ScanDelims(state.Pos, marker == 0x2A)

	for i := 1; i < 5; i++ {
		token := state.Push("text", "", 0)
		token.Content = string(marker)

		state.Delimiters = append(state.Delimiters, Delimiter{
			Marker: marker,
			Length: scanned.Length,
			Token:  len(state.Tokens) - 1,
			End:    -1,
			Open:   scanned.CanOpen,
			Close:  scanned.CanClose,
		})
	}

	state.Pos += scanned.Length

	return true
}

func (state *StateInline) _PostProcess(delimiters []Delimiter) {
	var startDelim Delimiter
	max := len(delimiters)

	for i := max; i >= 0; i-- {
		startDelim = delimiters[i]

		if startDelim.Marker != 0x5F /* _ */ && startDelim.Marker != 0x2A /* * */ {
			continue
		}

		// Process only opening markers
		if startDelim.End == -1 {
			continue
		}

		endDelim := delimiters[startDelim.End]

		// If the previous delimiter has the same marker and is adjacent to this one,
		// merge those into one strong delimiter.
		//
		// `<em><em>whatever</em></em>` -> `<strong>whatever</strong>`
		isStrong := i > 0 &&
			delimiters[i-1].End == startDelim.End+1 &&
			// check that first two markers match and adjacent
			delimiters[i-1].Marker == startDelim.Marker &&
			delimiters[i-1].Token == startDelim.Token-1 &&
			// check that last two markers are adjacent (we can safely assume they match)
			delimiters[startDelim.End+1].Token == endDelim.Token+1

		ch := string(startDelim.Marker)

		token := state.Tokens[startDelim.Token]

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

		token = state.Tokens[endDelim.Token]
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
			state.Tokens[delimiters[i-1].Token].Content = ""
			state.Tokens[delimiters[startDelim.End+1].Token].Content = ""
			i--
		}
	}
}

func Emphasis(
	_ *core.StateCore,
	_ *StateBlock,
	state *StateInline,
	_ int,
	_ int,
	silent bool,
) bool {
	return state.ETokenize(silent)
}

func EPostProcess(
	_ *core.StateCore,
	_ *StateBlock,
	state *StateInline,
	_ int,
	_ int,
	_ bool,
) bool {
	state.EPostProcess()
	return true
}

func (state *StateInline) EPostProcess() {
	tokensMeta := state.TokensMeta
	state._PostProcess(state.Delimiters)

	for _, tokenMeta := range tokensMeta {
		if len(tokenMeta.Delimiters) > 0 {
			state._PostProcess(tokenMeta.Delimiters)
		}
	}
}
