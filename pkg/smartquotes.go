package pkg

import (
	"unicode/utf8"
)

type StackValue struct {
	Token  int
	Pos    int
	Single bool
	Level  int
}
type Match struct {
	Index int
	Str   string
	Input string
}

func FindMatch(s string, start int, end int) (loc []Match) {
	text := Slice(s, start, end)

	match := QUOTE_TEST_RE.FindStringIndex(text)
	if len(match) == 0 {
		return nil
	}

	return []Match{
		{
			Index: match[0] + start,
			Str:   Slice(text, match[0], match[1]),
			Input: text,
		},
	}

}

func ProcessInline(tokens *[]*Token, state *StateCore) {
	var pos int
	var max int
	var i, j int
	var match Match
	var text string
	var token *Token
	var isSingle bool
	var thisLevel int
	var matches []Match
	var item StackValue
	var canOpen, canClose bool
	var lastChar, nextChar rune
	var openQuote, closeQuote string
	var isLastPunctChar, isNextPunctChar bool
	var isNextWhiteSpace, isLastWhiteSpace bool

	var stack []StackValue

	//fmt.Println("Attempting to process smartquotes")

	for i = 0; i < len(*tokens); i++ {
		token = (*tokens)[i]

		thisLevel = (*tokens)[i].Level

		for j = len(stack) - 1; j >= 0; j-- {
			if stack[j].Level <= thisLevel {
				break
			}
		}

		stack = stack[0 : j+1]

		if token.Type != "text" {
			continue
		}

		text = token.Content[0:]
		pos = 0
		max = utf8.RuneCountInString(text)

	OUTER:
		for pos < max {
			matches = FindMatch(text, pos, utf8.RuneCountInString(text))

			if matches == nil {
				break
			}

			match = matches[0]

			canOpen = true
			canClose = true
			pos = match.Index + 1
			isSingle = match.Str == "'"

			// Find previous character,
			// default to space if it's the beginning of the line
			lastChar = 0x20

			if match.Index-1 >= 0 {
				lastChar = CharCodeAt(text, match.Index-1)
			} else {
				for j = i - 1; j >= 0; j-- {
					if (*tokens)[j].Type == "softbreak" || (*tokens)[j].Type == "hardbreak" {
						break // lastChar defaults to 0x20
					}

					if utf8.RuneCountInString((*tokens)[j].Content) == 0 {
						continue // should skip all tokens except 'text', 'html_inline' or 'code_inline'
					}

					lastChar = CharCodeAt((*tokens)[j].Content, utf8.RuneCountInString((*tokens)[j].Content)-1)
					break
				}
			}

			// Find next character,
			// default to space if it's the end of the line
			//
			nextChar = 0x20

			if pos < max {
				nextChar = CharCodeAt(text, pos)
			} else {
				for j = i + 1; j < len(*tokens); j++ {
					if (*tokens)[j].Type == "softbreak" || (*tokens)[j].Type == "hardbreak" {
						//break; // nextChar defaults to 0x20
					}

					if utf8.RuneCountInString((*tokens)[j].Content) == 0 {
						continue // should skip all tokens except 'text', 'html_inline' or 'code_inline'
					}
					nextChar = CharCodeAt((*tokens)[j].Content, 0)
					break
				}
			}

			isLastPunctChar = IsMDAsciiPunct(lastChar) || IsPunctChar(lastChar)
			isNextPunctChar = IsMDAsciiPunct(nextChar) || IsPunctChar(nextChar)

			isLastWhiteSpace = IsWhiteSpace(lastChar)
			isNextWhiteSpace = IsWhiteSpace(nextChar)

			if isNextWhiteSpace {
				canOpen = false
			} else if isNextPunctChar {
				if !(isLastWhiteSpace || isLastPunctChar) {
					canOpen = false
				}
			}

			if isLastWhiteSpace {
				canClose = false
			} else if isLastPunctChar {
				if !(isNextWhiteSpace || isNextPunctChar) {
					canClose = false
				}
			}

			if nextChar == 0x22 /* " */ && match.Str == "\"" {
				if lastChar >= 0x30 /* 0 */ && lastChar <= 0x39 {
					// special case: 1"" - count first quote as an inch
					canClose = false
					canOpen = false
				}
			}

			if canOpen && canClose {
				// Replace quotes in the middle of punctuation sequence, but not
				// in the middle of the words, i.e.:
				//
				// 1. foo " bar " baz - not replaced
				// 2. foo-"-bar-"-baz - replaced
				// 3. foo"bar"baz     - not replaced
				//
				canOpen = isLastPunctChar
				canClose = isNextPunctChar
			}

			if !canOpen && !canClose {
				// middle of word
				if isSingle {
					token.Content = ReplaceAtIndex(token.Content, match.Index, "’")
				}
				continue
			}

			if canClose {
				// this could be a closing quote, rewind the stack to get a match
				for j = len(stack) - 1; j >= 0; j-- {
					item = stack[j]
					if stack[j].Level < thisLevel {
						break
					}
					if item.Single == isSingle && stack[j].Level == thisLevel {
						item = stack[j]

						if isSingle {
							openQuote = state.Md.Options.Quotes[2]
							closeQuote = state.Md.Options.Quotes[3]
						} else {
							openQuote = state.Md.Options.Quotes[0]
							closeQuote = state.Md.Options.Quotes[1]
						}

						// replace token.content *before* tokens[item.token].content,
						// because, if they are pointing at the same token, replaceAt
						// could mess up indices when quote length != 1
						token.Content = ReplaceAtIndex(token.Content, match.Index, closeQuote)
						(*tokens)[item.Token].Content = ReplaceAtIndex((*tokens)[item.Token].Content, item.Pos, openQuote)
						pos += utf8.RuneCountInString(closeQuote) - 1
						if item.Token == i {
							pos += utf8.RuneCountInString(openQuote) - 1
						}

						text = token.Content[:]
						max = utf8.RuneCountInString(text)

						stack = stack[0:j]
						continue OUTER
					}
				}
			}

			if canOpen {
				stack = append(stack, StackValue{
					Token:  i,
					Pos:    match.Index,
					Single: isSingle,
					Level:  thisLevel,
				})
			} else if canClose && isSingle {
				token.Content = ReplaceAtIndex(token.Content, match.Index, "’")
			}
		}
	}
}

func ReplaceAtIndex(s string, idx int, repl string) string {
	if utf8.RuneCountInString(s) == idx {
		return s + repl
	}

	runes := []rune(s)
	partOne := string(runes[:idx])
	partTwo := string(runes[idx+1:])
	return partOne + repl + partTwo
}

func Smartquotes(state *StateCore, _ *StateBlock, _ *StateInline, _ int, _ int, _ bool) bool {
	//fmt.Println("Processing Smartquotes")
	if !state.Md.Options.Typography {
		return false
	}

	for blkIdx := len(*state.Tokens) - 1; blkIdx >= 0; blkIdx-- {

		if (*state.Tokens)[blkIdx].Type != "inline" ||
			!QUOTE_TEST_RE.MatchString((*state.Tokens)[blkIdx].Content) {
			continue
		}

		ProcessInline(&(*state.Tokens)[blkIdx].Children, state)
	}
	return true
}
