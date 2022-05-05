package pkg

func BalancePairs(
	_ *StateCore,
	_ *StateBlock,
	state *StateInline,
	_ int,
	_ int,
	_ bool,
) bool {
	state.LinkPairs()
	return true
}

func (state *StateInline) ProcessDelimiters(delimiters []Delimiter) {

	max := len(delimiters)

	if max == 0 {
		return
	}

	headerIdx := 0
	var jumps []int
	var closerIdx int
	lastTokenIdx := -2
	var openersBottom map[rune][]int

	for closerIdx = 0; closerIdx < max; closerIdx++ {
		closer := delimiters[closerIdx]

		jumps = append(jumps, 0)

		// markers belong to same delimiter run if:
		//  - they have adjacent tokens
		//  - AND markers are the same

		if delimiters[headerIdx].Marker != closer.Marker ||
			lastTokenIdx != closer.Token-1 {
			headerIdx = closerIdx
		}

		lastTokenIdx = closer.Token

		// Length is only used for emphasis-specific "rule of 3",
		// if it's not defined (in strikethrough or 3rd party plugins),
		// we can default it to 0 to disable those checks.

		if closer.Length <= 1 {
			closer.Length = 0
		}

		if !closer.Close {
			continue
		}

		// Previously calculated lower bounds (previous fails)
		// for each marker, each delimiter length modulo 3,
		// and for whether this closer can be an opener;
		if _, ok := openersBottom[closer.Marker]; !ok {
			openersBottom[closer.Marker] = []int{-1, -1, -1, -1, -1, -1}
		}

		openIdx := 0

		if closer.Open {
			openIdx = 3
		}

		openIdx += closer.Length % 3

		minOpenerIdx := openersBottom[closer.Marker][openIdx]
		openerIdx := headerIdx - jumps[headerIdx] - 1
		newMinOpenerIdx := openerIdx

		for ; openerIdx > minOpenerIdx; openerIdx -= jumps[openerIdx] + 1 {
			opener := delimiters[openerIdx]

			if opener.Marker != closer.Marker {
				continue
			}

			if opener.Open && opener.End < 0 {
				isOddMatch := false

				// from spec:
				//
				// If one of the delimiters can both open and close emphasis, then the
				// sum of the lengths of the delimiter runs containing the opening and
				// closing delimiters must not be a multiple of 3 unless both lengths
				// are multiples of 3.
				//
				if opener.Close || closer.Open {
					if (opener.Length+closer.Length)%3 == 0 {
						if opener.Length%3 != 0 || closer.Length%3 != 0 {
							isOddMatch = true
						}
					}
				}

				if !isOddMatch {
					// If previous delimiter cannot be an opener, we can safely skip
					// the entire sequence in future checks. This is required to make
					// sure algorithm has linear complexity (see *_*_*_*_*_... case).

					lastJump := 0
					if openerIdx > 0 && !delimiters[openerIdx-1].Open {
						lastJump = jumps[openerIdx-1] + 1
					}

					jumps[closerIdx] = closerIdx - openerIdx + lastJump
					jumps[openerIdx] = lastJump

					closer.Open = false
					opener.End = closerIdx
					opener.Close = false
					newMinOpenerIdx = -1
					// treat next token as start of run,
					// it optimizes skips in **<...>**a**<...>** pathological case
					lastTokenIdx = -2
					break
				}
			}
		}
		if newMinOpenerIdx != -1 {
			// If match for this delimiter run failed, we want to set lower bound for
			// future lookups. This is required to make sure algorithm has linear
			// complexity.
			//
			// See details here:
			// https://github.com/commonmark/cmark/issues/178#issuecomment-270417442
			//

			openIdx = 0
			if closer.Open {
				openIdx = 3
			}

			if closer.Length > 0 {
				openIdx += closer.Length % 3
			}

			openersBottom[closer.Marker][openIdx] = newMinOpenerIdx
		}
	}
}

func (state *StateInline) LinkPairs() {
	tokensMeta := state.TokensMeta

	state.ProcessDelimiters(state.Delimiters)

	for _, tokenMeta := range tokensMeta {
		if len(tokenMeta.Delimiters) > 0 {
			state.ProcessDelimiters(tokenMeta.Delimiters)
		}
	}
}
