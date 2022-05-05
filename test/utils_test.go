package test

import (
	"github.com/stretchr/testify/assert"
	"go-markdown-it/pkg"
	"strings"

	"testing"
)

func TestFromCodePoint(t *testing.T) {
	assert.Equal(t, " ", pkg.FromCodePoint(0x20))
	assert.Equal(t, "😁", pkg.FromCodePoint(0x1F601))
}

func TestIsValidEntityCode(t *testing.T) {
	assert.Equal(t, true, pkg.IsValidEntityCode(0x20))

	assert.Equal(t, false, pkg.IsValidEntityCode(0x00))
	assert.Equal(t, false, pkg.IsValidEntityCode(0x0B))
	assert.Equal(t, false, pkg.IsValidEntityCode(0x0E))
	assert.Equal(t, false, pkg.IsValidEntityCode(0x7F))
	assert.Equal(t, false, pkg.IsValidEntityCode(0xD800))
	assert.Equal(t, false, pkg.IsValidEntityCode(0xFDD0))
	assert.Equal(t, false, pkg.IsValidEntityCode(0x1FFFF))
	assert.Equal(t, false, pkg.IsValidEntityCode(0x1FFFE))
}

func TestEscapeRegex(t *testing.T) {
	assert.Equal(t, " \\.\\?\\*\\+\\^\\$\\[\\]\\\\\\(\\)\\{\\}\\|\\-", pkg.EscapeRE(" .?*+^$[]\\(){}|-"))
}

func TestIsWhiteSpace(t *testing.T) {
	assert.Equal(t, true, pkg.IsWhiteSpace(0x09))
	assert.Equal(t, false, pkg.IsWhiteSpace(0x30))
	assert.Equal(t, true, pkg.IsWhiteSpace(0x2000))
}

func TestIsMdAsciiPunct(t *testing.T) {
	assert.Equal(t, false, pkg.IsMDAsciiPunct(0x30))

	chars := strings.Split("!\"#$%&\\'()*+,-./:;<=>?@[\\]^_`{|}~", "")

	for _, ch := range chars {
		assert.Equal(t, true, pkg.IsMDAsciiPunct(rune(ch[0])))
	}
}

func TestUnescapeMd(t *testing.T) {
	assert.Equal(t, "\\foo", pkg.UnescapedMD("\\foo"))

	chars := strings.Split("!\"#$%&\\'()*+,-./:;<=>?@[\\]^_`{|}~", "")

	for _, ch := range chars {
		assert.Equal(t, ch, pkg.UnescapedMD("\\"+ch))
	}
}
