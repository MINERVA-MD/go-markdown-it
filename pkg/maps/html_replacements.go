package maps

var HTML_REPLACEMENTS = map[string]string{
	"&":  "&amp;",
	"<":  "&lt;",
	">":  "&gt;",
	"\"": "&quot;",
}

var SCOPED_ABBR = map[string]string{
	"c":  "©",
	"r":  "®",
	"tm": "™",
}
