package pkg

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"regexp"
	"strings"
)

type HtmlSequence struct {
	Start     *regexp2.Regexp
	End       *regexp2.Regexp
	Terminate bool
}

var ATTR_NAME = "[a-zA-Z_:][a-zA-Z0-9:._-]*"
var UNQUOTED = "[^\"'=<>`\x00-\x20]+"
var SINGLE_QUOTED = "'[^']*'"
var DOUBLE_QUOTED = "\"[^\"]*\""
var ATTR_VALUE = fmt.Sprintf("(?:%s|%s|%s)", UNQUOTED, SINGLE_QUOTED, DOUBLE_QUOTED)
var ATTRIBUTE = fmt.Sprintf("(?:\\s+%s(?:\\s*=\\s*%s)?)", ATTR_NAME, ATTR_VALUE)
var OPEN_TAG = fmt.Sprintf("<[A-Za-z][A-Za-z0-9\\-]*%s*\\s*\\/?>", ATTRIBUTE)
var CLOSE_TAG = "<\\/[A-Za-z][A-Za-z0-9\\-]*\\s*>"
var COMMENT = "<!---->|<!--(?:-?[^>-])(?:-?[^-])*-->"
var PROCESSING = "<[?][\\s\\S]*?[?]>"
var DECLARATION = "<![A-Z]+\\s+[^>]*>"
var CDATA = "<!\\[CDATA\\[[\\s\\S]*?\\]\\]>"

var HTML_TAG_RE = regexp.MustCompile(fmt.Sprintf(
	"^(?:%s|%s|%s|%s|%s|%s)",
	OPEN_TAG,
	CLOSE_TAG,
	COMMENT,
	PROCESSING,
	DECLARATION,
	CDATA,
))

var HTML_OPEN_CLOSE_TAG_RE = regexp.MustCompile(fmt.Sprintf(
	"^(?:%s|%s)",
	OPEN_TAG,
	CLOSE_TAG,
))

var UNESCAPE_MD_RE = regexp.MustCompile("\\\\([!\"#$%&'()*+,\\-.\\/:;<=>?@[\\\\\\]^_`{|}~])")
var ENTITY_RE = regexp.MustCompile("(?i)&([a-z#][a-z0-9]{1,31});")
var UNESCAPE_ALL_RE = regexp.MustCompile(fmt.Sprintf(
	"(?i)%s|%s",
	UNESCAPE_MD_RE.String(),
	ENTITY_RE.String(),
))

var DIGITAL_ENTITY_TEST_RE = regexp.MustCompile("(?i)^#((?:x[a-f0-9]{1,8}|[0-9]{1,8}))")
var HTML_ESCAPE_TEST_RE = regexp.MustCompile("[&<>\"]")
var HTML_ESCAPE_REPLACE_RE = regexp.MustCompile("[&<>\"]")
var REGEXP_ESCAPE_RE = regexp.MustCompile("[.?*+^$[\\]\\\\(){}|-]")

//var UNICODE_PUNCT_RE = regexp.MustCompile(
//	"!-#%-\\*,-\\/:;\\?@\\[-\\]_\\{\\}\\xA1\\xA7\\xAB\\xB6\\xB7\\xBB\\xBF\\u037E\\u0387\\u055A-\\u055F\\u0589\\u058A\\u05BE\\u05C0\\u05C3\\u05C6\\u05F3\\u05F4\\u0609\\u060A\\u060C\\u060D\\u061B\\u061E\\u061F\\u066A-\\u066D\\u06D4\\u0700-\\u070D\\u07F7-\\u07F9\\u0830-\\u083E\\u085E\\u0964\\u0965\\u0970\\u09FD\\u0A76\\u0AF0\\u0C84\\u0DF4\\u0E4F\\u0E5A\\u0E5B\\u0F04-\\u0F12\\u0F14\\u0F3A-\\u0F3D\\u0F85\\u0FD0-\\u0FD4\\u0FD9\\u0FDA\\u104A-\\u104F\\u10FB\\u1360-\\u1368\\u1400\\u166D\\u166E\\u169B\\u169C\\u16EB-\\u16ED\\u1735\\u1736\\u17D4-\\u17D6\\u17D8-\\u17DA\\u1800-\\u180A\\u1944\\u1945\\u1A1E\\u1A1F\\u1AA0-\\u1AA6\\u1AA8-\\u1AAD\\u1B5A-\\u1B60\\u1BFC-\\u1BFF\\u1C3B-\\u1C3F\\u1C7E\\u1C7F\\u1CC0-\\u1CC7\\u1CD3\\u2010-\\u2027\\u2030-\\u2043\\u2045-\\u2051\\u2053-\\u205E\\u207D\\u207E\\u208D\\u208E\\u2308-\\u230B\\u2329\\u232A\\u2768-\\u2775\\u27C5\\u27C6\\u27E6-\\u27EF\\u2983-\\u2998\\u29D8-\\u29DB\\u29FC\\u29FD\\u2CF9-\\u2CFC\\u2CFE\\u2CFF\\u2D70\\u2E00-\\u2E2E\\u2E30-\\u2E4E\\u3001-\\u3003\\u3008-\\u3011\\u3014-\\u301F\\u3030\\u303D\\u30A0\\u30FB\\uA4FE\\uA4FF\\uA60D-\\uA60F\\uA673\\uA67E\\uA6F2-\\uA6F7\\uA874-\\uA877\\uA8CE\\uA8CF\\uA8F8-\\uA8FA\\uA8FC\\uA92E\\uA92F\\uA95F\\uA9C1-\\uA9CD\\uA9DE\\uA9DF\\uAA5C-\\uAA5F\\uAADE\\uAADF\\uAAF0\\uAAF1\\uABEB\\uFD3E\\uFD3F\\uFE10-\\uFE19\\uFE30-\\uFE52\\uFE54-\\uFE61\\uFE63\\uFE68\\uFE6A\\uFE6B\\uFF01-\\uFF03\\uFF05-\\uFF0A\\uFF0C-\\uFF0F\\uFF1A\\uFF1B\\uFF1F\\uFF20\\uFF3B-\\uFF3D\\uFF3F\\uFF5B\\uFF5D\\uFF5F-\\uFF65]|\\uD800[\\uDD00-\\uDD02\\uDF9F\\uDFD0]|\\uD801\\uDD6F|\\uD802[\\uDC57\\uDD1F\\uDD3F\\uDE50-\\uDE58\\uDE7F\\uDEF0-\\uDEF6\\uDF39-\\uDF3F\\uDF99-\\uDF9C]|\\uD803[\\uDF55-\\uDF59]|\\uD804[\\uDC47-\\uDC4D\\uDCBB\\uDCBC\\uDCBE-\\uDCC1\\uDD40-\\uDD43\\uDD74\\uDD75\\uDDC5-\\uDDC8\\uDDCD\\uDDDB\\uDDDD-\\uDDDF\\uDE38-\\uDE3D\\uDEA9]|\\uD805[\\uDC4B-\\uDC4F\\uDC5B\\uDC5D\\uDCC6\\uDDC1-\\uDDD7\\uDE41-\\uDE43\\uDE60-\\uDE6C\\uDF3C-\\uDF3E]|\\uD806[\\uDC3B\\uDE3F-\\uDE46\\uDE9A-\\uDE9C\\uDE9E-\\uDEA2]|\\uD807[\\uDC41-\\uDC45\\uDC70\\uDC71\\uDEF7\\uDEF8]|\\uD809[\\uDC70-\\uDC74]|\\uD81A[\\uDE6E\\uDE6F\\uDEF5\\uDF37-\\uDF3B\\uDF44]|\\uD81B[\\uDE97-\\uDE9A]|\\uD82F\\uDC9F|\\uD836[\\uDE87-\\uDE8B]|\\uD83A[\\uDD5E\\uDD5F]",
//)

var WHITESPACE_RE = regexp.MustCompile("\\s+")

var SPACE_RE = regexp.MustCompile("(\\s+)")
var LANG_ATTR = regexp.MustCompile("(\\s+)")

//(`(\s+)|\w`)

var NEWLINES_RE = regexp.MustCompile("\\r\\n?|\\n")

var HTML_SEQUENCES = []HtmlSequence{
	{
		Start:     regexp2.MustCompile("(?i)^<(script|pre|style|textarea)(?=(\\s|>|$))", 0),
		End:       regexp2.MustCompile("(?i)<\\/(script|pre|style|textarea)>", 0),
		Terminate: true,
	},
	{
		Start:     regexp2.MustCompile("^<!--", 0),
		End:       regexp2.MustCompile("-->", 0),
		Terminate: true,
	},
	{
		Start:     regexp2.MustCompile("^<\\?", 0),
		End:       regexp2.MustCompile("\\?>", 0),
		Terminate: true,
	},
	{
		Start:     regexp2.MustCompile("^<![A-Z]", 0),
		End:       regexp2.MustCompile(">", 0),
		Terminate: true,
	},
	{
		Start:     regexp2.MustCompile("^<!\\[CDATA\\[", 0),
		End:       regexp2.MustCompile("\\]\\]>", 0),
		Terminate: true,
	},
	{
		Start:     regexp2.MustCompile("(?i)"+"^</?("+strings.Join(HTML_BLOCKS[:], "|")+")(?=(\\s|/?>|$))", 0),
		End:       regexp2.MustCompile("^$", 0),
		Terminate: true,
	},
	{
		Start:     regexp2.MustCompile(HTML_OPEN_CLOSE_TAG_RE.String()+"\\s*$", 0),
		End:       regexp2.MustCompile("^$", 0),
		Terminate: false,
	},
}

var LINK_OPEN = regexp.MustCompile("(?i)^<a[>\\s]")
var LINK_CLOSE = regexp.MustCompile("(?i)^<\\/a\\s*>")

var DIGITAL_RE = regexp.MustCompile("(?i)^&#((?:x[a-f0-9]{1,6}|[0-9]{1,7}));")
var NAMED_RE = regexp.MustCompile("(?i)^&([a-z][a-z0-9]{1,31});")

var NEWLINE_RE = regexp.MustCompile("\n")
var BACKTICK_RE = regexp.MustCompile("^ (.+) $")

var EMAIL_RE = regexp.MustCompile("^([a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*)$")
var AUTOLINK_RE = regexp.MustCompile("^([a-zA-Z][a-zA-Z0-9+.\\-]{1,31}):([^<>\x00-\x20]*)$")
var TABLE_ALIGN_RE = regexp.MustCompile("^:?-+:?$")

var PLUS_MINUS_RE = regexp.MustCompile("\\+-")
var DOTS3_RE = regexp.MustCompile("\\.{2,}")
var QE_DOTS3_RE = regexp.MustCompile("([?!])???")
var QE4_RE = regexp.MustCompile("([?!]){4,}")
var COMMA_RE = regexp.MustCompile(",{2,}")
var EM_DASH_RE = regexp2.MustCompile("(?m)(^|[^-])---(?=[^-]|$)", 0)
var EN_DASH1_RE = regexp2.MustCompile("(?m)(^|\\s)--(?=\\s|$)", 0)
var EN_DASH2_RE = regexp2.MustCompile("(?m)(^|[^-\\s])--(?=[^-\\s]|$)", 0)
var RARE_RE = regexp.MustCompile("\\+-|\\.\\.|\\?\\?\\?\\?|!!!!|,,|--")
var SCOPED_ABBR_TEST_RE = regexp.MustCompile("(?i)\\((c|tm|r)\\)")
var SCOPED_ABBR_RE = regexp.MustCompile("(?i)\\((c|tm|r)\\)")

var QUOTE_TEST_RE = regexp.MustCompile("['\"]")

//var QUOTE_RE = regexp.MustCompile("['\"]")
//var APOSTROPHE = regexp.MustCompile("'\u2019")

var HTTP_RE = regexp.MustCompile("^http:\\/\\/")
var MAILTO_RE = regexp.MustCompile("(?i)^mailto:")

var SCHEME_RE = regexp.MustCompile("(?i)(?:^|[^a-z0-9.+-])([a-z][a-z0-9.+-]*)$")
var LINKIFY_CONFLICT_RE = regexp.MustCompile("\\*+$")

var BAD_PROTO_RE = regexp.MustCompile("^(vbscript|javascript|file|data):")
var GOOD_DATA_RE = regexp.MustCompile("^data:image\\/(gif|png|jpeg|webp);")

var NULL_RE = regexp.MustCompile("\u0000")
