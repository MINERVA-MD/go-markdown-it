package pkg

import (
	"unicode/utf8"
)

type TokenMeta struct {
	Delimiters *[]*Delimiter
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
	Src2             *MDString
	Md               *MarkdownIt
	Env              *Env
	Pos              int
	PosMax           int
	Pending          string
	Level            int
	PendingLevel     int
	Tokens           *[]*Token
	Cache            map[int]int
	Delimiters       *[]*Delimiter
	PrevDelimiters   []*[]*Delimiter
	Backticks        map[int]int
	BackTicksScanned bool
	LinkLevel        int
	TokensMeta       []TokenMeta
	//[]TokenMeta
}

func (state *StateInline) StateInline(src string, md *MarkdownIt, env *Env, outTokens *[]*Token) {
	mds := &MDString{}
	_ = mds.Init(src)

	state.Src = src
	state.Src2 = mds

	state.Env = env
	state.Md = md
	state.Tokens = outTokens
	state.TokensMeta = []TokenMeta{}

	state.Pos = 0
	state.PosMax = state.Src2.Length
	state.Level = 0
	state.Pending = ""

	state.PendingLevel = 0

	// Stores { start: end } pairs. Useful for backtrack
	// optimization of pairs parse (emphasis, strikes).
	state.Cache = map[int]int{}

	// List of emphasis-like delimiters for current tag
	state.Delimiters = &[]*Delimiter{}

	// Stack of delimiter lists for upper level tags
	state.PrevDelimiters = []*[]*Delimiter{}

	// backtick length => last seen position
	state.Backticks = map[int]int{}
	state.BackTicksScanned = false

	// Counter used to disable inline linkify-it execution
	// inside <a> and markdown links
	state.LinkLevel = 0
}

func (state *StateInline) Push(_type string, tag string, nesting int) *Token {

	if utf8.RuneCountInString(state.Pending) > 0 {
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

		oldDelimiters := Copy(*state.Delimiters)
		state.PrevDelimiters = append(state.PrevDelimiters, &oldDelimiters)

		state.Delimiters = &[]*Delimiter{}
		tokenMeta = TokenMeta{Delimiters: state.Delimiters}
		//TokenMeta{Delimiters: state.Delimiters}
	}

	state.PendingLevel = state.Level
	*state.Tokens = append(*state.Tokens, &token)

	state.TokensMeta = append(state.TokensMeta, tokenMeta)
	return &token
}

func Copy[T any](s []*T) []*T {
	c := make([]*T, len(s))
	for i, p := range s {

		if p == nil {
			// Skip to next for nil source pointer
			continue
		}

		// Create shallow copy of source element
		v := *p

		// Assign address of copy to destination.
		c[i] = &v
	}
	return c
}

func (state *StateInline) PushPending() *Token {
	token := GenerateToken("text", "", 0)

	token.Content = state.Pending
	token.Level = state.PendingLevel
	*state.Tokens = append(*state.Tokens, &token)
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
	marker, _ := state.Src2.CharCodeAt(start)

	// treat beginning of the line as a whitespace
	if start > 0 {
		lastChar, _ = state.Src2.CharCodeAt(start - 1)
	} else {
		lastChar = 0x20
	}

	for {
		if cc, _ := state.Src2.CharCodeAt(pos); pos < max && cc == marker {
			pos++
		} else {
			break
		}
	}

	var count = pos - start

	// treat end of the line as a whitespace
	if pos < max {
		nextChar, _ = state.Src2.CharCodeAt(pos)
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
