package pkg

type Helpers struct{}
type LinkResult struct {
	Ok    bool
	Pos   int
	Lines int
	Str   string
}

func (h *Helpers) ParseLinkLabel(state *StateInline, start int, disableNested bool) int {
	labelEnd := -1
	var found bool
	var marker rune
	var prevPos int
	max := state.PosMax
	oldPos := state.Pos

	state.Pos = start + 1
	level := 1

	for state.Pos < max {
		marker = CharCodeAt(state.Src, state.Pos)
		if marker == 0x5D /* ] */ {
			level--
			if level == 0 {
				found = true
				break
			}
		}

		prevPos = state.Pos
		state.Md.Inline.SkipToken(state)
		if marker == 0x5B /* [ */ {
			if prevPos == state.Pos-1 {
				// increase level if we find text `[`, which is not a part of any token
				level++
			} else if disableNested {
				state.Pos = oldPos
				return -1
			}
		}
	}

	if found {
		labelEnd = state.Pos
	}

	// restore old state
	state.Pos = oldPos

	return labelEnd
}

func (h *Helpers) ParseLinkDestination(str *MDString, pos int, max int) LinkResult {
	lines := 0
	start := pos
	var level int
	var code rune
	result := LinkResult{
		Ok:    false,
		Pos:   0,
		Lines: 0,
		Str:   "",
	}

	if cc, _ := str.CharCodeAt(pos); cc == 0x3C {
		pos++
		for pos < max {
			code, _ = str.CharCodeAt(pos)
			if code == 0x0A {
				return result
			}
			if code == 0x3C {
				return result
			}
			if code == 0x3E {
				result.Pos = pos + 1

				slice := str.Slice(start+1, pos)
				result.Str = UnescapeAll(slice)
				result.Ok = true
				return result
			}
			if code == 0x5C /* \ */ && pos+1 < max {
				pos += 2
				continue
			}

			pos++
		}

		// no closing '>'
		return result
	}

	// this should be ... } else { ... branch

	level = 0
	for pos < max {
		code, _ = str.CharCodeAt(pos)

		if code == 0x20 {
			break
		}

		// ascii control characters
		if code < 0x20 || code == 0x7F {
			break
		}

		if code == 0x5C /* \ */ && pos+1 < max {
			if cc, _ := str.CharCodeAt(pos + 1); cc == 0x20 {
				break
			}
			pos += 2
			continue
		}

		if code == 0x28 {
			level++
			if level > 32 {
				return result
			}
		}

		if code == 0x29 {
			if level == 0 {
				break
			}
			level--
		}

		pos++
	}

	if start == pos {
		return result
	}
	if level != 0 {
		return result
	}

	slice := str.Slice(start, pos)
	result.Str = UnescapeAll(slice)
	result.Lines = lines
	result.Pos = pos
	result.Ok = true

	return result
}

func (h *Helpers) ParseLinkTitle(str *MDString, pos int, max int) LinkResult {
	lines := 0
	start := pos
	var code rune
	var marker rune
	result := LinkResult{
		Ok:    false,
		Pos:   0,
		Lines: 0,
		Str:   "",
	}

	if pos >= max {
		return result
	}

	marker, _ = str.CharCodeAt(pos)

	if marker != 0x22 /* " */ && marker != 0x27 /* ' */ && marker != 0x28 {
		return result
	}

	pos++

	// if opening marker is "(", switch it to closing marker ")"
	if marker == 0x28 {
		marker = 0x29
	}

	for pos < max {
		code, _ = str.CharCodeAt(pos)
		if code == marker {
			result.Pos = pos + 1
			result.Lines = lines

			slice := str.Slice(start+1, pos)
			result.Str = UnescapeAll(slice)
			result.Ok = true
			return result
		} else if code == 0x28 /* ( */ && marker == 0x29 {
			return result
		} else if code == 0x0A {
			lines++
		} else if code == 0x5C /* \ */ && pos+1 < max {
			pos++
			if cc, _ := str.CharCodeAt(pos); cc == 0x0A {
				lines++
			}
		}
		pos++
	}

	return result
}
