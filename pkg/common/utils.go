package common

import (
	. "go-markdown-it/pkg/maps"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func IsValidEntityCode(c rune) bool {

	// broken sequence
	if c >= 0xD800 && c <= 0xDFFF {
		return false
	}
	// never used
	if c >= 0xFDD0 && c <= 0xFDEF {
		return false
	}
	if (c&0xFFFF) == 0xFFFF || (c&0xFFFF) == 0xFFFE {
		return false
	}
	// control codes
	if c >= 0x00 && c <= 0x08 {
		return false
	}
	if c == 0x0B {
		return false
	}
	if c >= 0x0E && c <= 0x1F {
		return false
	}
	if c >= 0x7F && c <= 0x9F {
		return false
	}

	// out of range
	if c > 0x10FFFF {
		return false
	}

	return true
}

func FromCodePoint(c rune) string {
	if c > 0xffff {
		c -= 0x10000
		var surrogate1 = 0xd800 + (c >> 10)
		var surrogate2 = 0xdc00 + (c & 0x3ff)
		return string([]rune{surrogate1, surrogate2})
	}
	return string(c)
}

func ReplaceEntityPattern(match string, name string) string {
	var code rune = 0

	if _, found := ENTITIES[name]; found {
		return ENTITIES[name]
	}

	var firstCharCode = []rune(name)[0]
	var isDigitalEntity = DIGITAL_ENTITY_TEST_RE.MatchString(name)

	if firstCharCode == 0x23 && isDigitalEntity {
		var secondChar = string([]rune(name)[1])

		if strings.ToLower(secondChar) == "x" {
			integer, _ := strconv.ParseInt(name[2:], 16, 0)
			code = rune(integer)
		} else {
			integer, _ := strconv.ParseInt(name[1:], 10, 0)
			code = rune(integer)
		}

		if IsValidEntityCode(code) {
			return FromCodePoint(code)
		}
	}

	return match
}

func UnescapedMD(str string) string {
	if !strings.Contains(str, "\\") {
		return str
	}
	return UNESCAPE_MD_RE.ReplaceAllString(str, "$1")
}

func UnescapedAll(str string) string {
	if !strings.Contains(str, "\\") &&
		!strings.Contains(str, "&") {
		return str
	}
	return ReplaceAllStringSubmatchFunc(UNESCAPE_ALL_RE, str,
		func(groups []string) string {
			var n = len(groups)
			var match = groups[0]

			if n > 1 {
				// escaped string
				return groups[1]
			}

			var entity = groups[2]
			return ReplaceEntityPattern(match, entity)

		})
}

func ReplaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func([]string) string) string {
	result := ""
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
		groups := []string{}
		for i := 0; i < len(v); i += 2 {
			if v[i] == -1 || v[i+1] == -1 {
				groups = append(groups, "")
			} else {
				groups = append(groups, str[v[i]:v[i+1]])
			}
		}

		result += str[lastIndex:v[0]] + repl(groups)
		lastIndex = v[1]
	}

	return result + str[lastIndex:]
}

func ReplaceUnsafChar(ch string) string {
	return HTML_REPLACEMENTS[ch]
}

func EscapeHtml(str string) string {
	var isHtmlEscape = HTML_ESCAPE_TEST_RE.MatchString(str)
	if isHtmlEscape {
		return HTML_ESCAPE_TEST_RE.ReplaceAllStringFunc(
			str,
			ReplaceUnsafChar,
		)
	}
	return str
}

func EscapeRE(str string) string {
	return REGEXP_ESCAPE_RE.ReplaceAllString(str, "\\$&")
}

func IsSpace(code rune) bool {

	if code == 0x09 || code == 0x20 {
		return true
	}

	return false
}

func IsWhiteSpace(code rune) bool {
	if code >= 0x2000 && code <= 0x200A {
		return true
	}

	switch code {
	case 0x09: // \t
		return true
	case 0x0A: // \n
		return true
	case 0x0B: // \v
		return true
	case 0x0C: // \f
		return true
	case 0x0D: // \r
		return true
	case 0x20:
		return true
	case 0xA0:
		return true
	case 0x1680:
		return true
	case 0x202F:
		return true
	case 0x205F:
		return true
	case 0x3000:
		return true

	default:
		return false
	}
}

func IsPunctChar(ch rune) bool {
	return unicode.IsPunct(ch)
}

func IsMDAsciiPunct(ch rune) bool {
	switch ch {
	case 0x21 /* ! */ :
		return true
	case 0x22 /* " */ :
		return true
	case 0x23 /* # */ :
		return true
	case 0x24 /* $ */ :
		return true
	case 0x25 /* % */ :
		return true
	case 0x26 /* & */ :
		return true
	case 0x27 /* ' */ :
		return true
	case 0x28 /* ( */ :
		return true
	case 0x29 /* ) */ :
		return true
	case 0x2A /* * */ :
		return true
	case 0x2B /* + */ :
		return true
	case 0x2C /* , */ :
		return true
	case 0x2D /* - */ :
		return true
	case 0x2E /* . */ :
		return true
	case 0x2F /* / */ :
		return true
	case 0x3A /* : */ :
		return true
	case 0x3B /* ; */ :
		return true
	case 0x3C /* < */ :
		return true
	case 0x3D /* = */ :
		return true
	case 0x3E /* > */ :
		return true
	case 0x3F /* ? */ :
		return true
	case 0x40 /* @ */ :
		return true
	case 0x5B /* [ */ :
		return true
	case 0x5C /* \ */ :
		return true
	case 0x5D /* ] */ :
		return true
	case 0x5E /* ^ */ :
		return true
	case 0x5F /* _ */ :
		return true
	case 0x60 /* ` */ :
		return true
	case 0x7B /* { */ :
		return true
	case 0x7C /* | */ :
		return true
	case 0x7D /* } */ :
		return true
	case 0x7E /* ~ */ :
		return true
	default:
		return false
	}
}

func NormalizeReference(str string) string {
	// Trim and collapse whitespace
	str = WHITESPACE_RE.ReplaceAllString(str, " ")
	return strings.ToUpper(strings.ToLower(str))
}
