package pkg

import "unicode/utf8"

func Newline(
	_ *StateCore,
	_ *StateBlock,
	state *StateInline,
	_ int,
	_ int,
	silent bool,
) bool {
	return state.Newline(silent)
}

func (state *StateInline) Newline(silent bool) bool {
	pos := state.Pos

	if cc, _ := state.Src2.CharCodeAt(pos); cc != 0x0A {
		return false
	}

	pmax := utf8.RuneCountInString(state.Pending) - 1
	max := state.PosMax

	// '  \n' -> hardbreak
	// Lookup in pending chars is bad practice! Don't copy to other rules!
	// Pending string is stored in concat mode, indexed lookups will cause
	// conversion to flat mode.

	if !silent {
		if pmax >= 0 && CharCodeAt(state.Pending, pmax) == 0x20 {
			if pmax >= 1 && CharCodeAt(state.Pending, pmax-1) == 0x20 {
				// Find whitespaces tail of pending chars.
				ws := pmax - 1

				// text_special
				for ws >= 1 && CharCodeAt(state.Pending, ws-1) == 0x20 {
					ws--
				}

				state.Pending = Slice(state.Pending, 0, ws)
				state.Push("hardbreak", "br", 0)
			} else {
				state.Pending = Slice(state.Pending, 0, utf8.RuneCountInString(state.Pending)-1)
				state.Push("softbreak", "br", 0)
			}
		} else {
			state.Push("softbreak", "br", 0)
		}
	}

	pos++

	for {
		if cc, _ := state.Src2.CharCodeAt(pos); pos < max && IsSpace(cc) {
			pos++
		} else {
			break
		}
	}

	state.Pos = pos

	return true
}
