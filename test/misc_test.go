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

	assert.Equal(t, "<p>123</p>\n", md.Render("123", pkg.Env{}))
	assert.Equal(t, "<p>&lt;!-- --&gt;</p>\n", md.Render("<!-- -->", pkg.Env{}))
}

func TestConfigureCoverage(t *testing.T) {

	var md = &pkg.MarkdownIt{}
	err := md.MarkdownIt("", pkg.Options{})

	assert.Equal(t, errors.New("wrong Markdown-It preset, can't be empty"), err)
}

func TestPlugin(t *testing.T) {
	// TODO
}

func TestCode(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{})
	assert.Equal(t, "<pre><code>hl\n</code></pre>\n", md.Render("```\nhl\n```", pkg.Env{}))
}

func TestCustomHighlight(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{Highlight: func(s string, s2 string, s3 string) string {
		return "<pre><code>==" + s + "==</code></pre>"
	}})

	assert.Equal(t, "<pre><code>==hl\n==</code></pre>\n", md.Render("```\nhl\n```", pkg.Env{}))
}
