package pkg

func IsTerminatorChar(ch rune) bool {
	switch ch {
	case 0x0A /* \n */ :
		return true
	case 0x21 /* ! */ :
		return true
	case 0x23 /* # */ :
		return true
	case 0x24 /* $ */ :
		return true
	case 0x25 /* % */ :
		return true
	case 0x26 /* & */ :
		return true
	case 0x2A /* * */ :
		return true
	case 0x2B /* + */ :
		return true
	case 0x2D /* - */ :
		return true
	case 0x3A /* : */ :
		return true
	case 0x3C /* < */ :
		return true
	case 0x3D /* = */ :
		return true
	case 0x3E /* > */ :
		return true
	case 0x40 /* @ */ :
		return true
	case 0x5B /* [ */ :
		return true
	case 0x5C /* \ */ :
		return true
	case 0x5D /* ] */ :
		return true
	case 0x5E /* ^ */ :
		return true
	case 0x5F /* _ */ :
		return true
	case 0x60 /* ` */ :
		return true
	case 0x7B /* { */ :
		return true
	case 0x7D /* } */ :
		return true
	case 0x7E /* ~ */ :
		return true
	default:
		return false
	}
}

func Text(
	_ *StateCore,
	_ *StateBlock,
	state *StateInline,
	_ int,
	_ int,
	silent bool,
) bool {
	//fmt.Println("Running Text")
	var pos = state.Pos

	for pos < state.PosMax && !IsTerminatorChar(CharCodeAt(state.Src, pos)) {
		pos++
	}

	if pos == state.Pos {
		return false
	}
	if !silent {
		state.Pending += string([]rune(state.Src)[state.Pos:pos])
	}

	state.Pos = pos

	return true
}
