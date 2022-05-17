package test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go-markdown-it/pkg"
	"testing"
	"unicode/utf16"
)

func TokenTypeFilter(tokens []*pkg.Token, _type string) []*pkg.Token {
	fTokens := []*pkg.Token{}

	for _, token := range tokens {
		if token.Type == _type {
			fTokens = append(fTokens, token)
		}
	}

	return fTokens
}

func TestMarkdownItConstructor(t *testing.T) {

	err := (&pkg.MarkdownIt{}).MarkdownIt("bad preset", pkg.Options{})
	assert.Equal(t, errors.New("wrong Markdown-It preset \"bad preset\", check name"), err)

	var md = &pkg.MarkdownIt{}
	err = md.MarkdownIt("commonmark", pkg.Options{Html: false})
	assert.Equal(t, nil, err)

	assert.Equal(t, "<p>123</p>\n", md.Render("123", &pkg.Env{}))
	assert.Equal(t, "<p>&lt;!-- --&gt;</p>\n", md.Render("<!-- -->", &pkg.Env{}))
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
	assert.Equal(t, "<pre><code>hl\n</code></pre>\n", md.Render("```\nhl\n```", &pkg.Env{}))
}

func TestCustomHighlight(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{Highlight: func(s string, _ string, _ string) string {
		return "<pre><code>==" + s + "==</code></pre>"
	}})

	assert.Equal(t, "<pre><code>==hl\n==</code></pre>\n", md.Render("```\nhl\n```", &pkg.Env{}))
}

func TestHighlightEscapeByDefault(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{Highlight: func(_ string, _ string, _ string) string {
		return ""
	}})

	assert.Equal(t, "<pre><code>&amp;\n</code></pre>\n", md.Render("```\n&\n```", &pkg.Env{}))
}

func TestHighlightArguments(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{Highlight: func(str string, lang string, attrs string) string {
		assert.Equal(t, "a", lang)
		assert.Equal(t, "b  c  d", attrs)
		return "<pre><code>==" + str + "==</code></pre>"
	}})

	assert.Equal(t, "<pre><code>==hl\n==</code></pre>\n", md.Render("``` a  b  c  d \nhl\n```", &pkg.Env{}))
}

func TestForceHardBreaks(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{Breaks: true})

	// TODO: Implement md.set() properly
	assert.Equal(t, "<p>a<br>\nb</p>\n", md.Render("a\nb", &pkg.Env{}))

	_ = md.MarkdownIt("default", pkg.Options{Breaks: true, XhtmlOut: true})
	assert.Equal(t, "<p>a<br />\nb</p>\n", md.Render("a\nb", &pkg.Env{}))
}

func TestXhtmlEnabled(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{XhtmlOut: true})

	assert.Equal(t, "<hr />\n", md.Render("---", &pkg.Env{}))
	assert.Equal(t, "<p><img src=\"\" alt=\"\" /></p>\n", md.Render("![]()", &pkg.Env{}))
	assert.Equal(t, "<p>a  <br />\nb</p>\n", md.Render("a  \\\nb", &pkg.Env{}))
}

func TestXhtmlDisabled(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{})

	assert.Equal(t, "<hr>\n", md.Render("---", &pkg.Env{}))
	assert.Equal(t, "<p><img src=\"\" alt=\"\"></p>\n", md.Render("![]()", &pkg.Env{}))
	assert.Equal(t, "<p>a  <br>\nb</p>\n", md.Render("a  \\\nb", &pkg.Env{}))
}

func TestEnableAndDisableRulesInChains(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{})

	was := map[string]int{
		"core":   len(md.Core.Ruler.GetRules("")),
		"block":  len(md.Block.Ruler.GetRules("")),
		"inline": len(md.Inline.Ruler.GetRules("")),
	}

	// Disable 2 rule in each chain & compare result
	_ = md.Disable([]string{"block", "inline", "code", "fence", "emphasis", "entity"}, false)

	now := map[string]int{
		"core":   len(md.Core.Ruler.GetRules("")) + 2,
		"block":  len(md.Block.Ruler.GetRules("")) + 2,
		"inline": len(md.Inline.Ruler.GetRules("")) + 2,
	}

	for k, v := range now {
		assert.Equal(t, v, was[k])
	}
}

func TestEnableAndDisableWithErrorControls(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{})

	err := md.Enable([]string{"link", "code", "invalid name"}, false)
	assert.Equal(t, errors.New("MarkdownIt. Failed to enable unknown rule(s): invalid name"), err)

	err = md.Disable([]string{"link", "code", "invalid name"}, false)
	assert.Equal(t, errors.New("MarkdownIt. Failed to disable unknown rule(s): invalid name"), err)

	err = md.Enable([]string{"link", "code"}, false)
	assert.Equal(t, nil, err)

	err = md.Disable([]string{"link", "code"}, false)
	assert.Equal(t, nil, err)
}

func TestShouldUnderstandStringsOnBulkEnableAndDisable(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{})

	_ = md.Disable([]string{"emphasis"}, false)
	assert.Equal(t, "_foo_", md.RenderInline("_foo_", &pkg.Env{}))

	//pretty
	_ = md.Enable([]string{"emphasis"}, false)
	assert.Equal(t, "<em>foo</em>", md.RenderInline("_foo_", &pkg.Env{}))
}

func TestShouldReplaceNullCharacters(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{})

	assert.Equal(t, "<p>foo\uFFFDbar</p>\n", md.Render("foo\u0000bar", &pkg.Env{}))
}

func TestShouldCorrectlyParseStringsWithTrailingNewline(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{})

	assert.Equal(t, "<p>123</p>\n", md.Render("123", &pkg.Env{}))
	assert.Equal(t, "<p>123</p>\n", md.Render("123\n", &pkg.Env{}))

	assert.Equal(t, "<pre><code>codeblock\n</code></pre>\n", md.Render("    codeblock", &pkg.Env{}))
	assert.Equal(t, "<pre><code>codeblock\n</code></pre>\n", md.Render("    codeblock\n", &pkg.Env{}))
}

func TestShouldQuicklyExitOnEmptyString(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{})

	assert.Equal(t, "", md.Render("", &pkg.Env{}))
}

func TestShouldParseInlinesOnly(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{})

	assert.Equal(t, "a <em>b</em> c", md.RenderInline("a *b* c", &pkg.Env{}))
}

// TODO: Implement pluggable renderer functions

func TestZeroPresetShouldDisableEverything(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("zero", pkg.Options{})

	assert.Equal(t, "<p>___foo___</p>\n", md.Render("___foo___", &pkg.Env{}))
	assert.Equal(t, "___foo___", md.RenderInline("___foo___", &pkg.Env{}))

	_ = md.Enable([]string{"emphasis"}, false)

	assert.Equal(t, "<p><em><strong>foo</strong></em></p>\n", md.Render("___foo___", &pkg.Env{}))
	assert.Equal(t, "<em><strong>foo</strong></em>", md.RenderInline("___foo___", &pkg.Env{}))
}

func TestShouldCheckBlockTerminationRulesWhenDisabled(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("zero", pkg.Options{})

	assert.Equal(t, "<p>foo\nbar</p>\n", md.Render("foo\nbar", &pkg.Env{}))
}

// TODO: Inline Plugin

func TestShouldNormalizeCRAndLF(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{})

	assert.Equal(t, md.Render("# test\r\r - hello\r - world\r", &pkg.Env{}), md.Render("# test\n\n - hello\n - world\n", &pkg.Env{}))
}

func TestShouldNormalizeCRAndLFToLF(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{})

	assert.Equal(t, md.Render("# test\r\n\r\n - hello\r\n - world\r\n", &pkg.Env{}), md.Render("# test\n\n - hello\n - world\n", &pkg.Env{}))
}

func TestShouldEscapeSurrogatePairs(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{})

	s := string(rune(0xD835))
	sp := string(utf16.DecodeRune(0xD835, 0xDC9C))

	assert.Equal(t, "<p>"+sp+"</p>\n", md.Render(sp, &pkg.Env{}))
	assert.Equal(t, "<p>"+s+"x</p>\n", md.Render(s+"x", &pkg.Env{}))
	assert.Equal(t, "<p>"+s+"</p>\n", md.Render(s, &pkg.Env{}))
}

// TODO: Url Normalization
// TODO: Link Validation

func TestBlockParserShouldNotNextAboveLimit(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{MaxNesting: 2})

	assert.Equal(t, "<blockquote>\n<p>foo</p>\n<blockquote></blockquote>\n</blockquote>\n", md.Render(">foo\n>>bar\n>>>baz", &pkg.Env{}))
}

func TestInlineParserShouldNotNextAboveLimit(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{MaxNesting: 1})

	assert.Equal(t, "<p><a href=\"\">`foo`</a></p>\n", md.Render("[`foo`]()", &pkg.Env{}))
}

func TestInlineNestingCoverage(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{MaxNesting: 2})

	assert.Equal(t, "<p>[[[[[[[[[[[[[[[[[[foo]()</p>\n", md.Render("[[[[[[[[[[[[[[[[[[foo]()", &pkg.Env{}))
}

func TestShouldSupportMultiCharQuotes(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{
		Typography: true,
		Quotes:     [4]string{"[[[", "]]", "(((((", "))))"},
	})

	assert.Equal(t, "<p>[[[foo]] (((((bar))))</p>\n", md.Render("\"foo\" 'bar'", &pkg.Env{}))
}

func TestShouldSupportNestedMultiCharQuotes(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{
		Typography: true,
		Quotes:     [4]string{"[[[", "]]", "(((((", "))))"},
	})

	assert.Equal(t, "<p>[[[foo (((((bar)))) baz]]</p>\n", md.Render("\"foo 'bar' baz\"", &pkg.Env{}))
}

func TestShouldSupportNestedMultiCharQuotesInDifferentTags(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{
		Typography: true,
		Quotes:     [4]string{"[[[", "]]", "(((((", "))))"},
	})

	assert.Equal(t, "<p>[[[a <em>b (((((c <em>d</em> e)))) f</em> g]]</p>\n", md.Render("\"a *b 'c *d* e' f* g\"", &pkg.Env{}))
}

func TestShouldMarkOrderedListItemTokensWithInfo(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{})

	tokens := md.Parse("1. Foo\n2. Bar\n20. Fuzz", &pkg.Env{})
	assert.Equal(t, 1, len(TokenTypeFilter(tokens, "ordered_list_open")))

	tokens = TokenTypeFilter(tokens, "list_item_open")
	assert.Equal(t, 3, len(tokens))
	assert.Equal(t, "1", tokens[0].Info)
	assert.Equal(t, ".", tokens[0].Markup)
	assert.Equal(t, "2", tokens[1].Info)
	assert.Equal(t, ".", tokens[1].Markup)
	assert.Equal(t, "20", tokens[2].Info)
	assert.Equal(t, ".", tokens[2].Markup)

	tokens = md.Parse(" 1. Foo\n2. Bar\n  20. Fuzz\n 199. Flp", &pkg.Env{})
	assert.Equal(t, 1, len(TokenTypeFilter(tokens, "ordered_list_open")))

	tokens = TokenTypeFilter(tokens, "list_item_open")
	assert.Equal(t, 4, len(tokens))
	assert.Equal(t, "1", tokens[0].Info)
	assert.Equal(t, ".", tokens[0].Markup)
	assert.Equal(t, "2", tokens[1].Info)
	assert.Equal(t, ".", tokens[1].Markup)
	assert.Equal(t, "20", tokens[2].Info)
	assert.Equal(t, ".", tokens[2].Markup)
	assert.Equal(t, "199", tokens[3].Info)
	assert.Equal(t, ".", tokens[3].Markup)
}

func TestShouldJoinTokenAttributes(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{})

	tokens := md.Parse("```", &pkg.Env{})

	tokens[0].AttrJoin(pkg.Attribute{
		Name:  "class",
		Value: "foo",
	})

	tokens[0].AttrJoin(pkg.Attribute{
		Name:  "class",
		Value: "bar",
	})

	assert.Equal(t, "<pre><code class=\"foo bar\"></code></pre>\n", md.Renderer.Render(tokens, md.Options, &pkg.Env{}))
}

func TestShouldSetTokenAttributes(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{})

	tokens := md.Parse("```", &pkg.Env{})

	tokens[0].AttrSet(pkg.Attribute{
		Name:  "class",
		Value: "foo",
	})

	assert.Equal(t, "<pre><code class=\"foo\"></code></pre>\n", md.Renderer.Render(tokens, md.Options, &pkg.Env{}))

	tokens[0].AttrSet(pkg.Attribute{
		Name:  "class",
		Value: "bar",
	})

	assert.Equal(t, "<pre><code class=\"bar\"></code></pre>\n", md.Renderer.Render(tokens, md.Options, &pkg.Env{}))
}

func TestShouldGetTokenAttributes(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("default", pkg.Options{})

	tokens := md.Parse("```", &pkg.Env{})

	attr, isPresent := tokens[0].AttrGet("myattr")
	assert.Equal(t, false, isPresent)

	tokens[0].AttrSet(pkg.Attribute{
		Name:  "myattr",
		Value: "myvalue",
	})

	attr, isPresent = tokens[0].AttrGet("myattr")
	assert.Equal(t, true, isPresent)
	assert.Equal(t, "myvalue", attr)
}
