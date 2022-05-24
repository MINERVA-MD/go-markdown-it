package pkg

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

	pmax := state.Pending2.Length - 1
	max := state.PosMax

	// '  \n' -> hardbreak
	// Lookup in pending chars is bad practice! Don't copy to other rules!
	// Pending string is stored in concat mode, indexed lookups will cause
	// conversion to flat mode.

	if !silent {
		if cc, _ := state.Pending2.CharCodeAt(pmax); pmax >= 0 && cc == 0x20 {
			if cc1, _ := state.Pending2.CharCodeAt(pmax - 1); pmax >= 1 && cc1 == 0x20 {
				// Find whitespaces tail of pending chars.
				ws := pmax - 1

				// text_special
				for {
					if ws1, _ := state.Pending2.CharCodeAt(ws - 1); ws >= 1 && ws1 == 0x20 {
						ws--
					} else {
						break
					}
				}

				_ = state.Pending2.Init(state.Pending2.Slice(0, ws))
				state.Push("hardbreak", "br", 0)
			} else {
				_ = state.Pending2.Init(state.Pending2.Slice(0, state.Pending2.Length-1))
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
