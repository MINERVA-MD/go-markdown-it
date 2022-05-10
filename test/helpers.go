package test

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type Options struct {
	Sep []string
}

type Content struct {
	Text  string
	Range []int
}

type Fixture struct {
	Type   string
	Header string
	First  Content
	Second Content
}

type Result struct {
	Meta     string
	Fixtures []Fixture
}

func FixLF(str string) string {
	if utf8.RuneCountInString(str) > 0 {
		return str + "\n"
	}
	return ""
}

func Parse(input string, options Options) Result {

	var SplitRe = regexp.MustCompile("\r?\n")
	var MetaRe = regexp.MustCompile("^-{3,}$")
	lines := SplitRe.Split(input, -1)

	min := 0
	line := 0

	var i int
	var l string
	max := len(lines)
	var result Result
	sep := options.Sep
	var blockStart int
	var currentSep string

	meta := ""

	if len(lines) > 0 &&
		utf8.RuneCountInString(strings.TrimSpace(lines[0])) > 0 {
		meta = strings.TrimSpace(lines[0])
	}

	if MetaRe.MatchString(meta) {
		line++
		for line < max && !MetaRe.MatchString(lines[line]) {
			line++
		}

		// If meta end found - extract range
		if line < max {
			result.Meta = strings.Join(lines[1:line], "\n")
			line++
			min = line

		} else {
			// if no meta closing - reset to start and try to parse data without meta
			line = 1
		}
	}

	// Scan fixtures
	for line < max {
		if !Contains(sep, lines[line]) {
			line++
			continue
		}

		currentSep = lines[line]

		fixture := Fixture{
			Type:   currentSep,
			Header: "",
			First: Content{
				Text:  "",
				Range: []int{},
			},
			Second: Content{
				Text:  "",
				Range: []int{},
			},
		}

		line++
		blockStart = line

		// seek end of first block
		for line < max && lines[line] != currentSep {
			line++
		}
		if line >= max {
			break
		}

		fixture.First.Text = FixLF(strings.Join(lines[blockStart:line], "\n"))
		fixture.First.Range = append(fixture.First.Range, blockStart, line)
		line++
		blockStart = line

		// seek end of second block
		for line < max && lines[line] != currentSep {
			line++
		}
		if line >= max {
			break
		}

		fixture.Second.Text = FixLF(strings.Join(lines[blockStart:line], "\n"))
		fixture.Second.Range = append(fixture.Second.Range, blockStart, line)
		line++

		// Look back for header on 2 lines before texture blocks
		i = fixture.First.Range[0] - 2
		for i >= Max(min, fixture.First.Range[0]-3) {
			l = lines[i]
			if Contains(sep, l) {
				break
			}
			if utf8.RuneCountInString(strings.TrimSpace(l)) > 0 {
				fixture.Header = strings.TrimSpace(l)
				break
			}
			i--
		}
		result.Fixtures = append(result.Fixtures, fixture)
	}
	return result
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
