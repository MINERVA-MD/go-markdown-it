package pkg

import (
	"bytes"
	"strings"
	"unicode/utf8"
)

type Rules struct{}
type Renderer struct {
	Rules Rules
}

// Running Text

func (r *Renderer) RenderAttrs(token *Token) string {
	var result = strings.Builder{}

	if token.Attrs == nil {
		return result.String()
	}

	for _, attr := range token.Attrs {
		result.WriteString(" ")
		result.WriteString(EscapeHtml(attr.Name))
		result.WriteString("=\"")
		result.WriteString(EscapeHtml(attr.Value))
		result.WriteString("\"")
	}

	return result.String()
}

func (r *Renderer) RenderToken(tokens []*Token, idx int, options Options) string {
	var nextToken *Token
	var result = strings.Builder{}
	var needLf = false
	var token = tokens[idx]

	// Tight list paragraphs
	if token.Hidden {
		return ""
	}

	if token.Block &&
		token.Nesting != -1 &&
		idx > 0 &&
		tokens[idx-1].Hidden {
		result.WriteString("\n")
	}

	// Add token name, e.g. `<img`
	if token.Nesting == -1 {
		result.WriteString("</")
		result.WriteString(token.Tag)
	} else {
		result.WriteString("<")
		result.WriteString(token.Tag)
	}

	// Encode attributes, e.g. `<img src="foo"`
	result.WriteString(r.RenderAttrs(token))

	// Add a slash for self-closing tags, e.g. `<img src="foo" /`
	if token.Nesting == 0 && options.XhtmlOut {
		result.WriteString(" /")
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
		result.WriteString(">\n")
	} else {
		result.WriteString(">")
	}

	return result.String()
}

func (r *Renderer) RenderInline(tokens []*Token, options Options, env *Env) string {
	var result = strings.Builder{}

	for idx, token := range tokens {
		if r.Rules.IsRuleTypeValid(token.Type) {
			result.WriteString(r.RenderRule(token.Type, tokens, idx, options, env))
		} else {
			result.WriteString(r.RenderToken(tokens, idx, options))
		}
	}
	return result.String()
}

func (r *Renderer) RenderInlineAsText(tokens []*Token, options Options, env Env) string {
	var result = strings.Builder{}

	for _, token := range tokens {
		if token.Type == "text" {
			result.WriteString(token.Content)
		} else if token.Type == "image" {
			result.WriteString(r.RenderInlineAsText(token.Children, options, env))
		} else if token.Type == "softbreak" {
			result.WriteString("\n")
		}
	}
	return result.String()
}

func (r *Renderer) Render(tokens []*Token, options Options, env *Env) string {
	var result = strings.Builder{}
	
	for idx, token := range tokens {
		if token.Type == "inline" {
			result.WriteString(r.RenderInline(token.Children, options, env))
		} else if r.Rules.IsRuleTypeValid(token.Type) {
			result.WriteString(r.RenderRule(token.Type, tokens, idx, options, env))
		} else {
			result.WriteString(r.RenderToken(tokens, idx, options))
		}
	}
	return result.String()
}

func (r *Renderer) RenderRule(Type string, tokens []*Token, idx int, options Options, env *Env) string {
	switch Type {
	case "code_inline":
		return r.Rules.CodeInline(tokens, idx, options, env, r)

	case "code_block":
		return r.Rules.CodeBlock(tokens, idx, options, env, r)

	case "fence":
		return r.Rules.Fence(tokens, idx, options, env, r)

	case "image":
		return r.Rules.Image(tokens, idx, options, *env, r)

	case "hardbreak":
		return r.Rules.Hardbreak(options)

	case "softbreak":
		return r.Rules.Softbreak(options)

	case "text":
		return r.Rules.Text(tokens, idx)

	case "html_block":
		return r.Rules.HtmlBlock(tokens, idx)

	case "html_inline":
		return r.Rules.HtmlInline(tokens, idx)

	default:
		return ""
	}
}

//===========================================================================

func (rules Rules) HtmlInline(tokens []*Token, idx int) string {
	return tokens[idx].Content
}

func (rules Rules) HtmlBlock(tokens []*Token, idx int) string {
	return tokens[idx].Content
}

func (rules Rules) Text(tokens []*Token, idx int) string {
	return EscapeHtml(tokens[idx].Content)
}

func (rules Rules) Softbreak(options Options) string {
	if options.Breaks {
		if options.XhtmlOut {
			return "<br />\n"
		}
		return "<br>\n"
	}
	return "\n"
}

func (rules Rules) Hardbreak(options Options) string {
	if options.XhtmlOut {
		return "<br />\n"
	}
	return "<br>\n"
}

func (rules Rules) Image(tokens []*Token, idx int, options Options, env Env, renderer *Renderer) string {
	var token = tokens[idx]
	var attrIdx = token.AttrIndex("alt")
	token.Attrs[attrIdx].Value = renderer.RenderInlineAsText(token.Children, options, env)

	return renderer.RenderToken(tokens, idx, options)
}

func (rules Rules) Fence(tokens []*Token, idx int, options Options, _ *Env, renderer *Renderer) string {

	var info string
	var arr []string
	var langName = ""
	var langAttrs = ""
	var highlighted = ""
	var tmpAttrs []Attribute

	var tmpToken *Token
	var token = tokens[idx]
	var result = strings.Builder{}

	if utf8.RuneCountInString(token.Info) > 0 {
		info = UnescapeAll(token.Info)
		info = strings.TrimSpace(info)
	} else {
		info = ""
	}

	if utf8.RuneCountInString(info) > 0 {
		//arr = LANG_ATTR.Split(info, -1)
		arr = SplitButIncludeDelimiter(info, " ")

		langName = arr[0]

		if len(arr) > 2 {
			langAttrs = strings.Join(arr[2:], "")
			langAttrs = strings.TrimSpace(langAttrs)
		} else {
			langAttrs = ""
		}

	}

	if options.Highlight != nil {
		var optHighlight = options.Highlight(token.Content, langName, langAttrs)
		if utf8.RuneCountInString(optHighlight) > 0 {
			highlighted = optHighlight
		} else {
			highlighted = EscapeHtml(token.Content)
		}
	} else {
		highlighted = EscapeHtml(token.Content)
	}

	if strings.Index(highlighted, "<pre") == 0 {
		return highlighted + "\n"
	}

	if utf8.RuneCountInString(info) > 0 {
		var i = token.AttrIndex("class")

		if token.Attrs != nil {
			copy(tmpAttrs, token.Attrs)
		} else {
			tmpAttrs = []Attribute{}
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
			tmpAttrs[i].Value += " " + options.LangPrefix + langName
		}

		tmpToken = &Token{Attrs: tmpAttrs}

		result.WriteString("<pre><code")
		result.WriteString(renderer.RenderAttrs(tmpToken))
		result.WriteString(">")
		result.WriteString(highlighted)
		result.WriteString("</code></pre>\n")

		return result.String()
	}

	result.WriteString("<pre><code")
	result.WriteString(renderer.RenderAttrs(token))
	result.WriteString(">")
	result.WriteString(highlighted)
	result.WriteString("</code></pre>\n")

	return result.String()
}

func (rules Rules) CodeBlock(tokens []*Token, idx int, _ Options, _ *Env, renderer *Renderer) string {
	var token = tokens[idx]
	var result = strings.Builder{}

	result.WriteString("<pre")
	result.WriteString(renderer.RenderAttrs(token))
	result.WriteString("><code>")
	result.WriteString(EscapeHtml(tokens[idx].Content))
	result.WriteString("</code></pre>\n")

	return result.String()
}

func (rules Rules) CodeInline(tokens []*Token, idx int, _ Options, _ *Env, renderer *Renderer) string {
	var token = tokens[idx]
	var result = strings.Builder{}

	result.WriteString("<code")
	result.WriteString(renderer.RenderAttrs(token))
	result.WriteString(">")
	result.WriteString(EscapeHtml(tokens[idx].Content))
	result.WriteString("</code>")

	return result.String()
}

func (rules Rules) IsRuleTypeValid(Type string) bool {

	switch Type {
	case "code_inline":
		return true

	case "code_block":
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

func SplitButIncludeDelimiter(s string, del string) []string {
	bs := []byte(s)
	sep := []byte(del)
	var ret [][]byte
	var splits []string

	for len(bs) > 0 {
		i := bytes.Index(bs, sep)
		if i == -1 {
			ret = append(ret, bs)
			break
		} else {
			ret = append(ret, bs[:i])
			ret = append(ret, sep)
			bs = bs[i+len(sep):]
		}
	}

	for _, split := range ret {
		splits = append(splits, string(split))
	}

	return splits
}
