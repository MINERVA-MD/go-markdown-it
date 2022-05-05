package test

import (
	"github.com/stretchr/testify/assert"
	"go-markdown-it/pkg"
	"testing"
)

func TestTokenAttr(t *testing.T) {
	token := pkg.GenerateToken("test_token", "tok", 1)

	assert.Equal(t, []pkg.Attribute(nil), token.Attrs)
	assert.Equal(t, -1, token.AttrIndex("foo"))

	token.AttrPush(pkg.Attribute{
		Name:  "foo",
		Value: "bar",
	})

	token.AttrPush(pkg.Attribute{
		Name:  "baz",
		Value: "bad",
	})

	assert.Equal(t, 0, token.AttrIndex("foo"))
	assert.Equal(t, 1, token.AttrIndex("baz"))
	assert.Equal(t, -1, token.AttrIndex("none"))
}
