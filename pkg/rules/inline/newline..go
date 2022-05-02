package inline

import (
	. "go-markdown-it/pkg/common"
	. "go-markdown-it/pkg/rules/block"
	"go-markdown-it/pkg/rules/core"
)

func Newline(
	_ *core.StateCore,
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

	if CharCodeAt(state.Src, pos) != 0x0A {
		return false
	}

	pmax := len(state.Pending) - 1
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

				for ws >= 1 && CharCodeAt(state.Pending, ws-1) == 0x20 {
					ws--
				}

				state.Pending = state.Pending[:ws]
				state.Push("softbreak", "br", 0)
			} else {
				state.Pending = state.Pending[:len(state.Pending)-1]
				state.Push("softbreak", "br", 0)
			}
		} else {
			state.Push("softbreak", "br", 0)
		}
	}

	pos++

	for pos < max && IsSpace(CharCodeAt(state.Src, pos)) {
		pos++
	}

	state.Pos = pos

	return true
}
