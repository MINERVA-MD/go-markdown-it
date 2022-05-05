package test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go-markdown-it/pkg"
	"testing"
)

func TestMarkdownItConstructor(t *testing.T) {

	err := (&pkg.MarkdownIt{}).MarkdownIt("bad preset", pkg.Options{})
	assert.Equal(t, errors.New("wrong Markdown-It preset \"bad preset\", check name"), err)

	var md = &pkg.MarkdownIt{}
	err = md.MarkdownIt("commonmark", pkg.Options{Html: true})
	assert.Equal(t, nil, err)

	assert.Equal(t, "<p>&lt;!-- --&gt;</p>\n", md.Render("<!-- -->", pkg.Env{}))
}
