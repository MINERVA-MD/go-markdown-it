package common

import (
	"fmt"
	"regexp"
)

var ATTR_NAME = "[a-zA-Z_:][a-zA-Z0-9:._-]*"
var UNQUOTED = "[^\"'=<>`\\x00-\\x20]+"
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
