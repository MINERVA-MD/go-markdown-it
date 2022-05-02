package inline

import (
	"go-markdown-it/pkg/common"
	"go-markdown-it/pkg/rules/block"
	"go-markdown-it/pkg/rules/core"
	"go-markdown-it/pkg/types"
)

func Linkify(
	_ *core.StateCore,
	_ *block.StateBlock,
	state *StateInline,
	_ int,
	_ int,
	silent bool,
) bool {
	return state.Linkify(silent)
}

func (state *StateInline) Linkify(silent bool) bool {

	if !state.Md.Options.Linkify {
		return false
	}

	if state.LinkLevel > 0 {
		return false
	}

	pos := state.Pos
	max := state.PosMax

	if pos+3 > max {
		return false
	}

	if block.CharCodeAt(state.Src, pos) != 0x3A {
		return false
	}
	if block.CharCodeAt(state.Src, pos+1) != 0x2F {
		return false
	}
	if block.CharCodeAt(state.Src, pos+2) != 0x2F {
		return false
	}

	match := common.SCHEME_RE.FindStringSubmatch(state.Pending)
	proto := match[1]

	// 	link := state.Md.Linkify.matchAtStart(state.src.slice(pos - proto.length));
	// TODO: Make proper call ^
	link := state.Src[pos-len(proto):]
	if len(link) == 0 {
		return false
	}

	// TODO: url = link.url
	url := link
	url = common.LINKIFY_CONFLICT_RE.ReplaceAllString(url, "")

	fullUrl := state.Md.NormalizeLink(url)
	if !state.Md.ValidateLink(fullUrl) {
		return false
	}

	if !silent {
		// TODO: double check negative start
		state.Pending = state.Pending[0:-len(proto)]

		token := state.Push("link_open", "a", 1);
		token.Attrs = []types.Attribute{
			{
				Name:  "href",
				Value: fullUrl,
			},
		}

		token.Markup = "linkify"
		token.Info = "auto"

		token = state.Push("text", "", 0)
		token.Content = state.Md.NormalizeLinkText(url)

		token = state.Push("link_close", "a", -1)
		token.Markup = "linkify"
		token.Info = "auto"
	}

	state.Pos += len(url) - len(proto)

	return true
}
