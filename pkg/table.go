package pkg

import (
	"strings"
	"unicode/utf8"
)

func (state *StateBlock) GetLine(line int) string {
	pos := state.BMarks[line] + state.TShift[line]
	max := state.EMarks[line]

	return state.Src2.Slice(pos, max)
}

func (state *StateBlock) EscapedSplit(str string) []string {
	var result []string
	pos := 0
	lastPos := 0
	var current = ""
	max := utf8.RuneCountInString(str)
	isEscaped := false

	ch := CharCodeAt(str, pos)

	for pos < max {
		if ch == 0x7c {
			if !isEscaped {
				// pipe separating cells, '|'

				result = append(result, current+Slice(str, lastPos, pos))
				current = ""
				lastPos = pos + 1
			} else {
				// escaped pipe, '\|'
				current += Slice(str, lastPos, pos-1)
				lastPos = pos
			}
		}

		isEscaped = ch == 0x5c /* \ */
		pos++

		ch = CharCodeAt(str, pos)
	}

	result = append(result, current+Slice(str, lastPos, utf8.RuneCountInString(str)))

	return result
}

func Table(
	_ *StateCore,
	state *StateBlock,
	_ *StateInline,
	startLine int,
	endLine int,
	silent bool,
) bool {
	return state.Table(startLine, endLine, silent)
}

func (state *StateBlock) Table(startLine int, endLine int, silent bool) bool {

	//fmt.Println("Processing Table")
	// should have at least two lines
	if startLine+2 > endLine {
		return false
	}

	nextLine := startLine + 1

	if state.SCount[nextLine] < state.BlkIndent {
		return false
	}

	// if it's indented more than 3 spaces, it should be a code block
	if state.SCount[nextLine]-state.BlkIndent >= 4 {
		return false
	}

	// first character of the second line should be '|', '-', ':',
	// and no other characters are allowed but spaces;
	// basically, this is the equivalent of /^[-:|][-:|\s]*$/ regexp

	pos := state.BMarks[nextLine] + state.TShift[nextLine]
	if pos >= state.EMarks[nextLine] {
		return false
	}

	firstCh, _ := state.Src2.CharCodeAt(pos)
	pos++

	if firstCh != 0x7C /* | */ && firstCh != 0x2D /* - */ && firstCh != 0x3A {
		return false
	}

	if pos >= state.EMarks[nextLine] {
		return false
	}

	secondCh, _ := state.Src2.CharCodeAt(pos)
	pos++

	if secondCh != 0x7C /* | */ && secondCh != 0x2D /* - */ && secondCh != 0x3A /* : */ && !IsSpace(secondCh) {
		return false
	}

	// if first character is '-', then second character must not be a space
	// (due to parsing ambiguity with list)
	if firstCh == 0x2D /* - */ && IsSpace(secondCh) {
		return false
	}

	for pos < state.EMarks[nextLine] {
		ch, _ := state.Src2.CharCodeAt(pos)

		if ch != 0x7C /* | */ && ch != 0x2D /* - */ && ch != 0x3A /* : */ && !IsSpace(ch) {
			return false
		}

		pos++
	}

	lineText := state.GetLine(startLine + 1)

	columns := strings.Split(lineText, "|")

	var aligns []string

	for i := 0; i < len(columns); i++ {
		t := strings.TrimSpace(columns[i])
		if utf8.RuneCountInString(t) == 0 {
			// allow empty columns before and after table, but not in between columns;
			// e.g. allow ` |---| `, disallow ` ---||--- `
			if i == 0 || i == len(columns)-1 {
				continue
			} else {
				return false
			}
		}

		if !TABLE_ALIGN_RE.MatchString(t) {
			return false
		}

		if CharCodeAt(t, utf8.RuneCountInString(t)-1) == 0x3A {
			if CharCodeAt(t, 0) == 0x3A {
				aligns = append(aligns, "center")
			} else {
				aligns = append(aligns, "right")
			}
		} else if CharCodeAt(t, 0) == 0x3A {
			aligns = append(aligns, "left")
		} else {
			aligns = append(aligns, "")
		}
	}

	lineText = strings.TrimSpace(state.GetLine(startLine))

	if strings.Index(lineText, "|") == -1 {
		return false
	}
	if state.SCount[startLine]-state.BlkIndent >= 4 {
		return false
	}
	columns = state.EscapedSplit(lineText)

	//fmt.Println(columns)

	if len(columns) > 0 && columns[0] == "" {
		columns = columns[1:]
	}

	if len(columns) > 0 && columns[len(columns)-1] == "" {
		columns = columns[:len(columns)-1]
	}

	// header row will define an amount of columns in the entire table,
	// and align row should be exactly the same (the rest of the rows can differ)
	columnCount := len(columns)
	if columnCount == 0 || columnCount != len(aligns) {
		return false
	}

	if silent {
		return true
	}

	oldParentType := state.ParentType
	state.ParentType = "table"

	// use 'blockquote' lists for termination because it's
	// the most similar to tables
	terminatorRules := state.Md.Block.Ruler.GetRules("blockquote")

	token := state.Push("table_open", "table", 1)
	tableLines := []int{startLine, 0}
	token.Map = []int{startLine, 0}

	token = state.Push("thead_open", "thead", 1)
	token.Map = []int{startLine, startLine + 1}

	token = state.Push("tr_open", "tr", 1)
	token.Map = []int{startLine, startLine + 1}

	for i := 0; i < len(columns); i++ {
		token = state.Push("th_open", "th", 1)
		if utf8.RuneCountInString(aligns[i]) > 0 {
			token.Attrs = []Attribute{
				{
					Name:  "style",
					Value: "text-align:" + aligns[i],
				},
			}
		}

		token = state.Push("inline", "", 0)
		token.Content = strings.TrimSpace(columns[i])
		token.Children = []*Token{}

		token = state.Push("th_close", "th", -1)
	}

	token = state.Push("tr_close", "tr", -1)
	token = state.Push("thead_close", "thead", -1)

	var tbodyLines []int
	for nextLine = startLine + 2; nextLine < endLine; nextLine++ {
		if state.SCount[nextLine] < state.BlkIndent {
			break
		}

		terminate := false
		l := len(terminatorRules)
		for i := 0; i < l; i++ {
			if terminatorRules[i](nil, state, nil, nextLine, endLine, true) {
				terminate = true
				break
			}
		}

		if terminate {
			break
		}
		lineText = strings.TrimSpace(state.GetLine(nextLine))
		if utf8.RuneCountInString(lineText) == 0 {
			break
		}
		if state.SCount[nextLine]-state.BlkIndent >= 4 {
			break
		}
		columns = state.EscapedSplit(lineText)
		if len(columns) > 0 && columns[0] == "" {
			columns = columns[1:]
		}
		if len(columns) > 0 && columns[len(columns)-1] == "" {
			columns = columns[:len(columns)-1]
		}

		if nextLine == startLine+2 {
			token = state.Push("tbody_open", "tbody", 1)
			tbodyLines = []int{startLine + 2, 0}
			token.Map = []int{startLine + 2, 0}
		}

		token = state.Push("tr_open", "tr", 1)
		token.Map = []int{nextLine, nextLine + 1}

		for i := 0; i < columnCount; i++ {
			token = state.Push("td_open", "td", 1)
			if i < len(aligns) && utf8.RuneCountInString(aligns[i]) > 0 {
				token.Attrs = []Attribute{
					{
						Name:  "style",
						Value: "text-align:" + aligns[i],
					},
				}
			}

			token = state.Push("inline", "", 0)

			n := len(columns)
			if i < n &&
				utf8.RuneCountInString(columns[i]) > 0 {
				token.Content = strings.TrimSpace(columns[i])
			} else {
				token.Content = ""
			}

			token.Children = []*Token{}

			token = state.Push("td_close", "td", -1)
		}
		token = state.Push("tr_close", "tr", -1)
	}

	if tbodyLines != nil && len(tbodyLines) > 0 {
		token = state.Push("tbody_close", "tbody", -1)
		tbodyLines[1] = nextLine
	}

	token = state.Push("table_close", "table", -1)
	tableLines[1] = nextLine

	state.ParentType = oldParentType
	state.Line = nextLine
	return true
}
