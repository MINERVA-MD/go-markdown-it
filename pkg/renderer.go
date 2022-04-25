package pkg

import (
	. "go-markdown-it/pkg/common"
	. "go-markdown-it/pkg/types"
	"strings"
)

func RenderAttrs(token *Token) string {
	var result = ""

	if token.Attrs == nil {
		return result
	}

	for _, attr := range token.Attrs {
		result += " " + EscapeHtml(attr.Name) + "=\"" + EscapeHtml(attr.Value)
	}

	return result
}

func RenderToken(tokens []*Token, idx int, options Options) string {
	var nextToken *Token
	var result = ""
	var needLf = false
	var token = tokens[idx]

	// Tight list paragraphs
	if token.Hidden {
		return result
	}

	if token.Block &&
		token.Nesting != -1 &&
		idx > 0 &&
		tokens[idx-1].Hidden {
		result += "\n"
	}

	// Add token name, e.g. `<img`
	if token.Nesting == -1 {
		result += "</" + token.Tag
	} else {
		result += "<" + token.Tag
	}

	// Encode attributes, e.g. `<img src="foo"`
	result += RenderAttrs(token)

	// Add a slash for self-closing tags, e.g. `<img src="foo" /`
	if token.Nesting == 0 && options.XhtmlOut {
		result += " /"
	}

	// Check if we need to add a newline after this tag
	if token.Block {
		needLf = true
		if token.Nesting == 1 {
			if idx+1 < len(tokens) {
				nextToken = tokens[idx+1]
				if nextToken.Type == "inline" || nextToken.Hidden {
					// Block-level tag containing an inline tag.
					needLf = false
				} else if nextToken.Nesting == -1 && nextToken.Tag == token.Tag {
					// Opening tag + closing tag of the same type. E.g. `<li></li>`.
					needLf = false
				}
			}
		}
	}

	if needLf {
		result += ">\n"
	} else {
		result += ">"
	}

	return result
}

func RenderInline(tokens []*Token, options Options, env Env) string {
	var result string

	for idx, token := range tokens {
		if IsRuleTypeValid(token.Type) {
			result += RenderRule(token.Type, tokens, idx, options, env)
		} else {
			result += RenderToken(tokens, idx, options)
		}
	}
	return result
}

func RenderInlineAsText(tokens []*Token, options Options, env Env) string {
	var result = ""

	for _, token := range tokens {
		if token.Type == "text" {
			result += token.Content
		} else if token.Type == "image" {
			result += RenderInlineAsText(token.Children, options, env)
		} else if token.Type == "softbreak" {
			result += "\n"
		}
	}
	return result
}

func Render(tokens []*Token, options Options, env Env) string {
	var result = ""

	for idx, token := range tokens {
		if token.Type == "inline" {
			result += RenderInline(token.Children, options, env)
		} else if IsRuleTypeValid(token.Type) {
			result += RenderRule(token.Type, tokens, idx, options, env)
		} else {
			result += RenderToken(tokens, idx, options)
		}
	}
	return result
}

func RenderRule(Type string, tokens []*Token, idx int, options Options, env Env) string {
	switch Type {
	case "code_inline":
		return CodeInline(tokens, idx, options, env)

	case "code_block":
		return CodeBlock(tokens, idx, options, env)

	case "fence":
		return Fence(tokens, idx, options, env)

	case "image":
		return Image(tokens, idx, options, env)

	case "hardbreak":
		return Hardbreak(options)

	case "softbreak":
		return Softbreak(options)

	case "text":
		return Text(tokens, idx)

	case "html_block":
		return HtmlBlock(tokens, idx)

	case "html_inline":
		return HtmlInline(tokens, idx)

	default:
		return ""
	}
}

//===========================================================================

func HtmlInline(tokens []*Token, idx int) string {
	return tokens[idx].Content
}

func HtmlBlock(tokens []*Token, idx int) string {
	return tokens[idx].Content
}

func Text(tokens []*Token, idx int) string {
	return EscapeHtml(tokens[idx].Content)
}

func Softbreak(options Options) string {
	if options.Breaks {
		if options.XhtmlOut {
			return "<br />\n"
		}
		return "<br>\n"
	}
	return "\n"
}

func Hardbreak(options Options) string {
	if options.XhtmlOut {
		return "<br />\n"
	}
	return "<br>\n"
}

func Image(tokens []*Token, idx int, options Options, env Env) string {
	var token = tokens[idx]
	var attrIdx = AttrIndex(token, "alt")
	token.Attrs[attrIdx].Value = RenderInlineAsText(tokens, options, env)

	return RenderToken(tokens, idx, options)
}

func Fence(tokens []*Token, idx int, options Options, env Env) string {

	var info string
	var arr []string
	var langName = ""
	var langAttrs = ""
	var highlighted = ""
	var tmpAttrs []Attribute

	var tmpToken *Token
	var token = tokens[idx]

	if len(token.Info) > 0 {
		info = strings.TrimSpace(token.Info)
	} else {
		info = ""
	}

	if len(info) > 0 {
		arr = SPACE_RE.Split(info, 2)
		langName = arr[0]
		langAttrs = arr[1]
	}

	if options.Highlight != nil {
		var optHighlight = options.Highlight(token.Content, langName, langAttrs)
		if len(optHighlight) > 0 {
			highlighted = optHighlight
		} else {
			highlighted = EscapeHtml(token.Content)
		}
	}

	if strings.Index(highlighted, "<pre") == 0 {
		return highlighted + "\n"
	}

	if len(info) > 0 {
		var i = AttrIndex(token, "class")
		if token.Attrs != nil {
			copy(tmpAttrs, token.Attrs)
		} else {
			token.Attrs = []Attribute{}
		}

		if i < 0 {
			tmpAttrs = append(tmpAttrs, Attribute{
				Name:  "class",
				Value: options.LangPrefix + langName,
			})
		} else {
			// This call makes no sense: Requires testing
			tmpAttrs[i] = Attribute{
				Name:  tmpAttrs[i].Name,
				Value: tmpAttrs[i].Value,
			}
		}
		tmpAttrs[i].Value += " " + options.LangPrefix + langName

		tmpToken = &Token{Attrs: tmpAttrs}

		return "<pre><code" + RenderAttrs(tmpToken) + ">" +
			highlighted +
			"</code></pre>\n"
	}

	return "<pre><code" + RenderAttrs(token) + ">" +
		highlighted +
		"</code></pre>\n"
}

func CodeBlock(tokens []*Token, idx int, options Options, env Env) string {
	var token = tokens[idx]

	return "<pre" + RenderAttrs(token) + "><code>" +
		EscapeHtml(tokens[idx].Content) +
		"</code></pre>\n"
}

func CodeInline(tokens []*Token, idx int, options Options, env Env) string {
	var token = tokens[idx]

	return "<code" + RenderAttrs(token) + ">" +
		EscapeHtml(tokens[idx].Content) +
		"</code>"
}

func IsRuleTypeValid(Type string) bool {

	switch Type {
	case "code_inline":
		return true

	case "code_block ":
		return true

	case "fence":
		return true

	case "image":
		return true

	case "hardbreak":
		return true

	case "softbreak":
		return true

	case "text":
		return true

	case "html_block":
		return true

	case "html_inline":
		return true

	default:
		return false

	}
}
