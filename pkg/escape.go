package pkg

var ESCAPED = [256]int{}

func InitEscapedChars() {
	chars := "\\!\"#$%&'()*+,./:;<=>?@[]^_`{|}~-"

	for _, char := range chars {
		ESCAPED[char] = 1
	}
}

func Escape(
	_ *StateCore,
	_ *StateBlock,
	state *StateInline,
	_ int,
	_ int,
	silent bool,
) bool {

	//fmt.Println("Entered Escape")
	pos := state.Pos
	max := state.PosMax
	initEscapeMap := false

	if CharCodeAt(state.Src, pos) != 0x5C {
		return false
	}
	pos++

	if pos >= max {
		return false
	}

	ch1 := CharCodeAt(state.Src, pos)
	var ch2 rune

	if ch1 == 0x0A {
		if !silent {
			state.Push("hardbreak", "br", 0)
		}

		pos++
		// skip leading whitespaces from next line

		for pos < max {
			ch1 = CharCodeAt(state.Src, pos)
			if !IsSpace(ch1) {
				break
			}
			pos++
		}

		state.Pos = pos
		return true
	}

	// TODO: Double check this indexing is Unicode compliant
	escapedStr := string(CharCodeAt(state.Src, pos))

	if ch1 >= 0xD800 && ch1 <= 0xDBFF && pos+1 < max {
		ch2 = CharCodeAt(state.Src, pos+1)

		if ch2 >= 0xDC00 && ch2 <= 0xDFFF {
			// TODO: Double check this indexing is Unicode compliant
			escapedStr += string(state.Src[pos+1])
			pos++
		}
	}

	origStr := `\` + escapedStr

	if !silent {
		//fmt.Println("Got here~~~~~~~~~~~")
		token := state.Push("text_special", "", 0)

		//fmt.Println(ch1)
		if !initEscapeMap {
			InitEscapedChars()
			initEscapeMap = true
		}

		if ch1 < 256 && ESCAPED[ch1] != 0 {
			token.Content = escapedStr
		} else {
			token.Content = origStr
		}

		token.Markup = origStr
		token.Info = "escape"
	}

	state.Pos = pos + 1
	return true
}
