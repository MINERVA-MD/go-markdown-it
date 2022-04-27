package inline

import (
	. "go-markdown-it/pkg"
	. "go-markdown-it/pkg/common"
	. "go-markdown-it/pkg/rules/block"
	. "go-markdown-it/pkg/types"
)

type TokenMeta struct {
	Delimiters []Delimiter
}

type Delimiter struct {
	Marker rune
	Length int
	Token  int
	End    int
	Open   bool
	Close  bool
}

type DelimScan struct {
	CanOpen  bool
	CanClose bool
	Length   int
}

type StateInline struct {
	Src              string
	Md               Parser
	Env              Env
	Pos              int
	PosMax           int
	Pending          string
	Level            int
	PendingLevel     int
	Tokens           []*Token
	Cache            map[string]string
	Delimiters       []Delimiter
	PrevDelimiters   [][]Delimiter
	Backticks        string
	BackTicksScanned bool
	LinkLevel        int
	TokensMeta       []TokenMeta
}

func (state *StateInline) Push(_type string, tag string, nesting int) *Token {
	if len(state.Pending) > 0 {
		state.PushPending()
	}

	token := GenerateToken(_type, tag, nesting)

	if nesting < 0 {
		state.Level--
		state.Delimiters, state.PrevDelimiters = Pop(state.PrevDelimiters)
	}

	token.Level = state.Level
	var tokenMeta TokenMeta

	if nesting > 0 {
		// Opening Tag
		state.Level++
		state.PrevDelimiters = append(state.PrevDelimiters, state.Delimiters)
		state.Delimiters = []Delimiter{}
		tokenMeta = TokenMeta{Delimiters: state.Delimiters}
	}

	state.PendingLevel = state.Level
	state.Tokens = append(state.Tokens, &token)
	state.TokensMeta = append(state.TokensMeta, tokenMeta)

	return &token
}

func (state *StateInline) PushPending() *Token {
	token := GenerateToken("text", "", 0)

	token.Content = state.Pending
	token.Level = state.PendingLevel
	state.Tokens = append(state.Tokens, &token)
	state.Pending = ""

	return &token
}

func (state *StateInline) ScanDelims(start int, canSplitWord bool) DelimScan {

	pos := start
	var lastChar rune
	var nextChar rune
	max := state.PosMax
	leftFlanking := true
	rightFlanking := true
	marker := CharCodeAt(state.Src, start)

	// treat beginning of the line as a whitespace
	if start > 0 {
		lastChar = CharCodeAt(state.Src, start-1)
	} else {
		lastChar = 0x20
	}

	for pos < max && CharCodeAt(state.Src, pos) == marker {
		pos++
	}

	var count = pos - start

	// treat end of the line as a whitespace
	if pos < max {
		nextChar = CharCodeAt(state.Src, pos)
	} else {
		nextChar = 0x20
	}

	isLastPunctChar := IsMDAsciiPunct(lastChar) || IsPunctChar(lastChar)
	isNextPunctChar := IsMDAsciiPunct(nextChar) || IsPunctChar(nextChar)

	isLastWhiteSpace := IsWhiteSpace(lastChar)
	isNextWhiteSpace := IsWhiteSpace(nextChar)

	if isNextWhiteSpace {
		leftFlanking = false
	} else if isNextPunctChar {
		if !(isLastWhiteSpace || isLastPunctChar) {
			leftFlanking = false
		}
	}

	if isLastWhiteSpace {
		rightFlanking = false
	} else if isLastPunctChar {
		if !(isNextWhiteSpace || isNextPunctChar) {
			rightFlanking = false
		}
	}

	var canOpen bool
	var canClose bool

	if !canSplitWord {
		canOpen = leftFlanking && (!rightFlanking || isLastPunctChar)
		canClose = rightFlanking && (!leftFlanking || isNextPunctChar)
	} else {
		canOpen = leftFlanking
		canClose = rightFlanking
	}

	return DelimScan{
		CanOpen:  canOpen,
		CanClose: canClose,
		Length:   count,
	}
}

func Pop[T any](s []T) (T, []T) {
	lastElem := s[len(s)-1]
	collection := s[:len(s)-1]

	return lastElem, collection
}
