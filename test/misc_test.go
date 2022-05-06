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
	_ = md.MarkdownIt("default", pkg.Options{Highlight: func(s string, _ string, _ string) string {
		return "<pre><code>==" + s + "==</code></pre>"
	}})

	assert.Equal(t, "<pre><code>==hl\n==</code></pre>\n", md.Render("```\nhl\n```", pkg.Env{}))
}

func TestHighlightEscapeByDefault(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{Highlight: func(_ string, _ string, _ string) string {
		return ""
	}})

	assert.Equal(t, "<pre><code>&amp;\n</code></pre>\n", md.Render("```\n&\n```", pkg.Env{}))
}

func TestHighlightArguments(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{Highlight: func(str string, lang string, attrs string) string {
		assert.Equal(t, "a", lang)
		assert.Equal(t, "b  c  d", attrs)
		return "<pre><code>==" + str + "==</code></pre>"
	}})

	assert.Equal(t, "<pre><code>==hl\n==</code></pre>\n", md.Render("``` a  b  c  d \nhl\n```", pkg.Env{}))
}

func TestForceHardBreaks(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{Breaks: true})

	// TODO: Implement md.set() properly
	assert.Equal(t, "<p>a<br>\nb</p>\n", md.Render("a\nb", pkg.Env{}))

	_ = md.MarkdownIt("default", pkg.Options{Breaks: true, XhtmlOut: true})
	assert.Equal(t, "<p>a<br />\nb</p>\n", md.Render("a\nb", pkg.Env{}))
}
