package pkg

import (
	"strings"
	"unicode/utf8"
)

type StateBlock struct {
	Src        string
	Src2       *MDString
	Md         *MarkdownIt
	Env        *Env
	Tokens     *[]*Token
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

func (state *StateBlock) StateBlock(src string, md *MarkdownIt, env *Env, outTokens *[]*Token) {
	mds := &MDString{}
	_ = mds.Init(src)

	state.Src = src
	state.Src2 = mds

	state.Md = md
	state.Env = env
	state.Tokens = outTokens

	state.BMarks = []int{}
	state.EMarks = []int{}
	state.TShift = []int{}
	state.SCount = []int{}

	state.BsCount = []int{}

	// block parser variables
	state.BlkIndent = 0
	state.Line = 0
	state.LineMax = 0
	state.Tight = false
	state.DDIndent = -1
	state.ListIndent = -1

	state.ParentType = "root"

	state.Level = 0
	state.Result = ""

	// Create caches
	// Generate markers.

	var pos = 0
	var start = 0
	var indent = 0
	var offset = 0

	var s = state.Src2
	var n = s.Length

	var indentFound = false

	for ; pos < n; pos++ {

		ch, _ := s.CharCodeAt(pos)

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

	*state.Tokens = append(*state.Tokens, &token)
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
	var max = state.Src2.Length
	for ; pos < max; pos++ {
		ch, _ := state.Src2.CharCodeAt(pos)
		if !IsSpace(ch) {
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
		ch, _ := state.Src2.CharCodeAt(pos)

		if !IsSpace(ch) {
			return pos + 1
		}
	}

	return pos
}

func (state *StateBlock) SkipChars(pos int, code rune) int {
	var max = state.Src2.Length
	for ; pos < max; pos++ {
		ch, _ := state.Src2.CharCodeAt(pos)
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
		ch, _ := state.Src2.CharCodeAt(pos)

		if rune(code) != ch {
			return pos + 1
		}
	}

	return pos
}

func (state *StateBlock) GetLines(begin int, end int, indent int, keepLastLF bool) string {

	var ch rune
	var last int
	var first int
	var line = begin
	var lineStart int
	var lineIndent int

	if begin >= end {
		return ""
	}

	queue := make([]string, end-begin)

	for i := 0; line < end; i++ {
		lineIndent = 0
		first = state.BMarks[line]
		lineStart = first

		if line+1 < end || keepLastLF {
			// No need for bounds check because we have fake entry on tail.
			last = state.EMarks[line] + 1
		} else {
			last = state.EMarks[line]
		}

		for first < last && lineIndent < indent {
			ch, _ = state.Src2.CharCodeAt(first)

			if IsSpace(ch) {
				if ch == 0x09 {
					lineIndent += 4 - (lineIndent+state.BsCount[line])%4
				} else {
					lineIndent++
				}
			} else if first-lineStart < state.TShift[line] {
				// patched tShift masked characters to look like spaces (blockquotes, list markers)
				lineIndent++
			} else {
				break
			}

			first++
		}

		slice := state.Src2.Slice(first, last)
		if lineIndent > indent {
			// partially expanding tabs in code blocks, e.g '\t\tfoobar'
			// with indent=2 becomes '  \tfoobar'

			queue[i] = strings.Join(make([]string, lineIndent-indent+1)[:], " ") + slice
		} else {
			queue[i] = slice
		}
		line++
	}
	return strings.Join(queue[:], "")
}

// CharCodeAt This is O(n), consider replacing this for optimization’s sake
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

func Slice(s string, start int, end int) string {
	if end <= start {
		return ""
	}

	n := utf8.RuneCountInString(s)
	if start >= n || end > n {
		return ""
	}
	return string([]rune(s)[start:end])
}
