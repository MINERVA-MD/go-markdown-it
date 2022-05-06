package pkg

import (
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	"strings"
)

//type Config struct {
//	Default    types.Preset
//	Zero       types.Preset
//	Commonmark types.Preset
//}

var config = map[string]Preset{
	"default":    DefaultPresets,
	"zero":       ZeroPresets,
	"commonmark": CommonmarkPresets,
}

////////////////////////////////////////////////////////////////////////////////
//
// This validator can prohibit more than really needed to prevent XSS. It's a
// tradeoff to keep code simple and to be secure by default.
//
// If you need different setup - override validator method as you wish. Or
// replace it with dummy function and use external sanitizer.

func ValidateLink(url string) bool {
	// url should be normalized at this point, and existing entities are decoded
	str := strings.TrimSpace(url)
	str = strings.ToLower(str)

	if BAD_PROTO_RE.MatchString(str) {
		if GOOD_DATA_RE.MatchString(str) {
			return true
		} else {
			return false
		}
	}
	return true
}

////////////////////////////////////////////////////////////////////////////////

func NormalizeLink(url string) string {
	// TODO
	return url
}

func NormalizeLinkText(url string) string {
	// TODO
	return url
}

type MarkdownIt struct {
	Inline            ParserInline
	Block             ParserBlock
	Core              ParserCore
	Renderer          Renderer
	Linkify           Linkify
	ValidateLink      func(string) bool
	NormalizeLink     func(string) string
	NormalizeLinkText func(string) string
	Options           Options
	Helpers           Helpers
}

/**
 * Main parser/renderer class.
 *
 * ##### Usage
 *
 * ```javascript
 * // node.js, "classic" way:
 * var MarkdownIt = require('markdown-it'),
 *     md = new MarkdownIt();
 * var result = md.render('# markdown-it rulezz!');
 *
 * // node.js, the same, but with sugar:
 * var md = require('markdown-it')();
 * var result = md.render('# markdown-it rulezz!');
 *
 * // browser without AMD, added to "window" on script load
 * // Note, there are no dash.
 * var md = window.markdownit();
 * var result = md.render('# markdown-it rulezz!');
 * ```
 *
 * Single line rendering, without paragraph wrap:
 *
 * ```javascript
 * var md = require('markdown-it')();
 * var result = md.renderInline('__markdown-it__ rulezz!');
 * ```
 **/

// MarkdownIt - Main parser/renderer class.
func (md *MarkdownIt) MarkdownIt(presetName string, options Options) error {

	md.Inline = ParserInline{Ruler: Ruler{}}
	md.Block = ParserBlock{Ruler: Ruler{}}
	md.Core = ParserCore{Ruler: Ruler{}}
	md.Renderer = Renderer{Rules: Rules{}}

	md.Inline.ParserInline()
	md.Block.ParserBlock()
	md.Core.ParserCore()

	// TODO: Not attached to correct struct
	md.Linkify = Linkify{}

	md.ValidateLink = ValidateLink
	md.NormalizeLink = NormalizeLink
	md.NormalizeLinkText = NormalizeLinkText

	// TODO: Handle error or let it propagate
	err := md.Configure(presetName)

	if err != nil {
		return err
	}

	// No conditional needed since we'll
	// at least pass in some default settings
	md.Set(options)

	if md.Options.MaxNesting == 0 {
		md.Options.MaxNesting = 100
	}

	return nil
}

func (md *MarkdownIt) Configure(presetName string) error {

	if len(presetName) == 0 {
		return errors.New("wrong Markdown-It preset, can't be empty")

	}

	if _, ok := config[presetName]; !ok {
		return errors.New(fmt.Sprintf("wrong Markdown-It preset \"%s\", check name", presetName))
	}

	presets := config[presetName]

	// ParserCore
	if len(presets.Components.Core.Rules) > 0 {
		md.Core.Ruler.EnableOnly(presets.Components.Core.Rules, false)
	}

	// Block
	if len(presets.Components.Block.Rules) > 0 {
		md.Block.Ruler.EnableOnly(presets.Components.Block.Rules, false)
	}

	// Inline
	if len(presets.Components.Inline.Rules) > 0 {
		md.Inline.Ruler.EnableOnly(presets.Components.Inline.Rules, false)
	}

	// Inline (2)
	if len(presets.Components.Inline.Rules2) > 0 {
		md.Inline.Ruler2.EnableOnly(presets.Components.Inline.Rules2, false)
	}

	return nil
}

func (md *MarkdownIt) Set(options Options) {
	md.Options = options
}

func (md *MarkdownIt) Enable(list []string, ignoreInvalid bool) error {

	var _list []string
	var result []string
	for _, ruler := range []Ruler{md.Core.Ruler, md.Block.Ruler, md.Block.Ruler} {
		_list, _ = ruler.Enable(list, true)
		result = append(result, _list...)
	}

	// Ruler 2
	_list, _ = md.Inline.Ruler2.Enable(list, true)
	result = append(result, _list...)

	var missed []string

	for _, name := range list {
		if !slices.Contains(result, name) {
			missed = append(missed, name)
		}
	}

	if len(missed) == 0 && !ignoreInvalid {
		return errors.New("MarkdownIt. Failed to enable unknown rule(s): " + missed[0])
	}

	return nil
}

func (md *MarkdownIt) Disable(list []string, ignoreInvalid bool) error {

	var _list []string
	var result []string
	for _, ruler := range []Ruler{md.Core.Ruler, md.Block.Ruler, md.Block.Ruler} {
		_list, _ = ruler.Disable(list, true)
		result = append(result, _list...)
	}

	// Ruler 2
	_list, _ = md.Inline.Ruler2.Disable(list, true)
	result = append(result, _list...)

	var missed []string

	for _, name := range list {
		if !slices.Contains(result, name) {
			missed = append(missed, name)
		}
	}

	if len(missed) == 0 && !ignoreInvalid {
		return errors.New("MarkdownIt. Failed to disable unknown rule(s): " + missed[0])
	}

	return nil
}

func (md *MarkdownIt) Use(_ string) {

}

func (md *MarkdownIt) Parse(src string, env Env) []*Token {
	var state = &StateCore{}
	state.StateCore(src, md, env)

	md.Core.Process(state)
	//utils.PrettyPrint(state.Tokens)

	return state.Tokens
}

func (md *MarkdownIt) Render(src string, env Env) string {
	tokens := md.Parse(src, env)
	return md.Renderer.Render(tokens, md.Options, env)
}

func (md *MarkdownIt) ParseInline(src string, env Env) []*Token {
	var state = &StateCore{}
	state.StateCore(src, md, env)
	state.InlineMode = true

	md.Core.Process(state)

	return state.Tokens
}

func (md *MarkdownIt) RenderInline(src string, env Env) string {
	tokens := md.ParseInline(src, env)
	return md.Renderer.Render(tokens, md.Options, env)
}
