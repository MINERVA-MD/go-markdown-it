package block

import (
	. "go-markdown-it/pkg"
	. "go-markdown-it/pkg/common"
	. "go-markdown-it/pkg/types"
)

type StateBlock struct {
	Src        string
	Md         Parser
	Env        Env
	Tokens     []*Token
	BMarks     []int
	EMarks     []int
	TShift     []int
	SCount     []int
	BsCount    []int
	BlkIndent  int
	Line       int
	LineMax    int
	Tight      bool
	DDIndent   int
	ListIndent int
	ParentType string
	Level      int
	Result     string
}

func StateBlockInit() *StateBlock {
	state := StateBlock{
		Src:        "",
		Md:         Parser{},
		Env:        Env{},
		Tokens:     nil,
		BMarks:     nil,
		EMarks:     nil,
		TShift:     nil,
		SCount:     nil,
		BsCount:    nil,
		BlkIndent:  0,
		Line:       0,
		LineMax:    0,
		Tight:      false,
		DDIndent:   0,
		ListIndent: 0,
		ParentType: "",
		Level:      0,
		Result:     "",
	}

	return &state
}

func (state *StateBlock) StateBlock() {
	// TODO
	var s = state.Src
	var start = 0
	var pos = 0
	var indent = 0
	var offset = 0
	var n = len(s)
	var indentFound = false

	for ; pos < n; pos++ {
		ch := CharCodeAt(s, pos)

		if !indentFound {
			if IsSpace(ch) {
				indent++

				if ch == 0x09 {
					offset += 4 - offset%4
				} else {
					offset++
				}
				continue
			} else {
				indentFound = true
			}
		}

		if ch == 0x0A || pos == n-1 {
			if ch != 0x0A {
				pos++
			}
			state.BMarks = append(state.BMarks, start)
			state.EMarks = append(state.EMarks, pos)
			state.TShift = append(state.TShift, indent)
			state.SCount = append(state.SCount, offset)
			state.BsCount = append(state.BsCount, 0)

			indentFound = false
			indent = 0
			offset = 0
			start = pos + 1
		}
	}

	// Push fake entry to simplify cache bounds checks
	state.BMarks = append(state.BMarks, n)
	state.EMarks = append(state.EMarks, n)
	state.TShift = append(state.TShift, 0)
	state.SCount = append(state.SCount, 0)
	state.BsCount = append(state.BsCount, 0)

	state.LineMax = len(state.BMarks) - 1
}

func (state *StateBlock) Push(_type string, tag string, nesting int) *Token {
	var token = GenerateToken(_type, tag, nesting)
	token.Block = true

	if nesting < 0 { // closing tag
		state.Level--
	}
	token.Level = state.Level
	if nesting > 0 { // opening tag
		state.Level++
	}

	state.Tokens = append(state.Tokens, &token)
	return &token
}

func (state *StateBlock) IsEmpty(line int) bool {
	return state.BMarks[line]+state.TShift[line] >= state.EMarks[line]
}

func (state *StateBlock) SkipEmptyLines(from int) int {
	var max = state.LineMax

	for ; from < max; from++ {
		if state.BMarks[from]+state.TShift[from] < state.EMarks[from] {
			break
		}
	}

	return from
}

func (state *StateBlock) SkipSpaces(pos int) int {
	var max = len(state.Src)
	for ; pos < max; pos++ {
		ch := CharCodeAt(state.Src, pos)
		if IsSpace(ch) {
			break
		}
	}
	return pos
}

func (state *StateBlock) SkipSpacesBack(pos int, min int) int {
	if pos <= min {
		return pos
	}

	for pos > min {
		pos--
		ch := CharCodeAt(state.Src, pos)

		if !IsSpace(ch) {
			return pos + 1
		}
	}

	return pos
}

func (state *StateBlock) SkipChars(pos int, code rune) int {
	var max = len(state.Src)
	for ; pos < max; pos++ {
		ch := CharCodeAt(state.Src, pos)
		if ch != code {
			break
		}
	}
	return pos
}

func (state *StateBlock) SkipCharsBack(pos int, code int, min int) int {
	if pos <= min {
		return pos
	}

	for pos > min {
		pos--
		ch := CharCodeAt(state.Src, pos)

		if rune(code) != ch {
			return pos + 1
		}
	}

	return pos
}

func (state *StateBlock) GetLines(begin int, end int, indent int, keepLastLF bool) string {
	// TODO
	return ""
}

// CharCodeAt This is O(n) consider replacing this for optimization sake
func CharCodeAt(s string, n int) rune {
	i := 0
	for _, r := range s {
		if i == n {
			return r
		}
		i++
	}
	return 0
}
