package common

import (
	"fmt"
	. "go-markdown-it/pkg/maps"
)

func IsValidEntityCode(c rune) bool {

	// broken sequence
	if c >= 0xD800 && c <= 0xDFFF {
		return false
	}
	// never used
	if c >= 0xFDD0 && c <= 0xFDEF {
		return false
	}
	if (c&0xFFFF) == 0xFFFF || (c&0xFFFF) == 0xFFFE {
		return false
	}
	// control codes
	if c >= 0x00 && c <= 0x08 {
		return false
	}
	if c == 0x0B {
		return false
	}
	if c >= 0x0E && c <= 0x1F {
		return false
	}
	if c >= 0x7F && c <= 0x9F {
		return false
	}

	// out of range
	if c > 0x10FFFF {
		return false
	}

	return true
}

func FromCodePoint(c rune) string {
	if c > 0xffff {
		c -= 0x10000
		var surrogate1 = 0xd800 + (c >> 10)
		var surrogate2 = 0xdc00 + (c & 0x3ff)
		return string([]rune{surrogate1, surrogate2})
	}
	return string(c)
}

func ReplaceEntityPattern(match string, name string) {
	var code rune = 0

	if v, found := ENTITIES[name]; found {
		fmt.Println(v)
	}
}
