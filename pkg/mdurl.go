package pkg

import (
	"github.com/dlclark/regexp2"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var encodeCache []string
var componentChars = "-_.!~*'()"
var defaultChars = ";/?:@&=+$,-_.!~*'()#"
var twoHexRe = regexp.MustCompile("(?i)^[0-9a-f]{2}")
var alphaCharsRe = regexp.MustCompile("(?i)^[0-9a-z]$")
var hostRe = regexp.MustCompile("^//[^@/]+@[^@/]+")
var protocolPatternRe = regexp.MustCompile("(?i)^([a-z0-9.+-]+:)")
var portPatternRe = regexp.MustCompile(":[0-9]*$")
var simplePathPatternRe = regexp2.MustCompile("^(//?(?!/)[^\\?\\s]*)(\\?[^\\s]*)?$", 0)
var delims = []string{"<", ">", "\"", "`", " ", "\r", "\n", "\t"}
var unwise = append([]string{"{", "}", "|", "\\", "^", "`"}, delims[:]...)
var autoEscape = append([]string{"'"}, unwise[:]...)
var nonHostChars = append([]string{"%", "/", "?", ";", "#"}, autoEscape[:]...)
var hostEndingChars = []string{"/", "?", "#"}
var hostnameMaxLen = 255
var hostnamePartPatternRe = regexp.MustCompile("^[+a-z0-9A-Z_-]{0,63}$")
var hostnamePartStart = regexp.MustCompile("^([+a-z0-9A-Z_-]{0,63})(.*)$")
var hostlessProtocol = map[string]bool{
	"javascript":  true,
	"javascript:": true,
}

var slashedProtocol = map[string]bool{
	"http":    true,
	"https":   true,
	"ftp":     true,
	"gopher":  true,
	"file":    true,
	"http:":   true,
	"https:":  true,
	"ftp:":    true,
	"gopher:": true,
	"file:":   true,
}

type Url struct {
	Protocol string
	Slashes  bool
	Auth     string
	Port     string
	Hostname string
	Hash     string
	Search   string
	Pathname string
}
type MdUrl struct{}

func (md *MdUrl) Parse(url string, slashesDenoteHost bool) Url {
	rest := url
	retUrl := Url{}
	var hec int
	var slashes bool
	var proto_ string
	var lowerProto string
	rest = strings.TrimSpace(rest)

	if !slashesDenoteHost && len(strings.Split(url, "#")) == 1 {
		// Try fast path regexp
		if m, _ := simplePathPatternRe.FindStringMatch(rest); m != nil {
			// you can get all the groups too
			gps := m.Groups()

			retUrl.Pathname = gps[1].Captures[0].String()
			if len(gps) > 1 &&
				utf8.RuneCountInString(gps[2].Captures[0].String()) > 0 {
				retUrl.Search = gps[2].Captures[0].String()
			}
			return retUrl
		}
	}

	var proto = protocolPatternRe.FindAllString(rest, -1)
	if len(proto) > 0 {
		proto_ = proto[0]
		lowerProto = strings.ToLower(proto_)
		retUrl.Protocol = proto_
		rest = string(([]rune(rest))[utf8.RuneCountInString(proto_):])
	}

	// figure out if it's got a host
	// user@server is *always* interpreted as a hostname, and url
	// resolution will treat //foo/bar as host=foo,path=bar because that's
	// how the browser resolves relative URLs.
	if slashesDenoteHost || utf8.RuneCountInString(proto_) > 0 || hostRe.MatchString(rest) {
		slashes = string(([]rune(rest))[0:2]) == "//"

		if slashes &&
			!(utf8.RuneCountInString(proto_) > 0 && hostlessProtocol[proto_]) {
			rest = string(([]rune(rest))[2:])
			retUrl.Slashes = true
		}
		//rest = rest.substr(2);
		//this.slashes = true;
		//}
	}

	if !hostlessProtocol[proto_] &&
		(slashes || (utf8.RuneCountInString(proto_) > 0 && !slashedProtocol[proto_])) {

		// there's a hostname.
		// the first instance of /, ?, ;, or # ends the host.
		//
		// If there is an @ in the hostname, then non-host chars *are* allowed
		// to the left of the last @ sign, unless some host-ending character
		// comes *before* the @-sign.
		// URLs are obnoxious.
		//
		// ex:
		// http://a@b@c/ => user:a@b host:c
		// http://a@b?@c => user:a host:c path:/?@c

		// v0.12 TODO(isaacs): This is not quite how Chrome does things.
		// Review our test case against browsers more comprehensively.

		// find the first instance of any hostEndingChars
		var hostEnd = -1
		for i := 0; i < len(hostEndingChars); i++ {
			hec = strings.Index(rest, hostEndingChars[i])
			if hec != -1 && (hostEnd == -1 || hec < hostEnd) {
				hostEnd = hec
			}
		}

		// at this point, either we have an explicit point where the
		// auth portion cannot go past, or the last @ char is the decider.
		var atSign int
		var auth string

		if hostEnd == -1 {
			// atSign can be anywhere.
			atSign = strings.LastIndex(rest, "@")
		} else {
			// atSign must be in auth portion.
			// http://a@b/c@d => host:b auth:a path:/c@d
			at := strings.LastIndex(string([]rune(rest)[hostEnd:]), "")

			if at < 0 {
				atSign = -1
			} else {
				atSign = hostEnd + at
			}
		}

		// Now we have a portion which is definitely the auth.
		// Pull that off.
		if atSign != -1 {
			auth = string(([]rune(rest))[0:atSign])
			rest = string(([]rune(rest))[atSign+1:])

			retUrl.Auth = auth
		}

		// the host is the remaining to the left of the first non-host char
		hostEnd = -1
		for i := 0; i < len(nonHostChars); i++ {
			hec = strings.Index(rest, nonHostChars[i])
			if hec != -1 && (hostEnd == -1 || hec < hostEnd) {
				hostEnd = hec
			}
		}
		//// if we still have not hit it, then the entire thing is a host.
		if hostEnd == -1 {
			hostEnd = utf8.RuneCountInString(rest)
		}
		//
		if string([]rune(rest)[hostEnd-1]) == ":" {
			hostEnd--
		}

		var host = string(([]rune(rest))[0:hostEnd])
		rest = string(([]rune(rest))[hostEnd:])

		//// pull out port.
		md.ParseHost(host, &retUrl)

		//// we've indicated that there is a hostname,
		//// so even if it's empty, it has to be present.

		//retUrl.Hostname = this.hostname || ''
		//

		//// if hostname begins with [ and ends with ]
		//// assume that it's an IPv6 address.
		var ipv6Hostname = string([]rune(retUrl.Hostname)[0]) == "[" &&
			string([]rune(retUrl.Hostname)[utf8.RuneCountInString(retUrl.Hostname)-1]) == "]"

		//// validate a little.
		if !ipv6Hostname {
			var hostparts = strings.Split(retUrl.Hostname, ".")
			l := len(hostparts)
			for i := 0; i < l; i++ {
				var part = hostparts[i]
				if utf8.RuneCountInString(part) == 0 {
					continue
				}
				if !hostnamePartPatternRe.MatchString(part) {
					var newpart = ""
					k := utf8.RuneCountInString(part)
					for j := 0; j < k; j++ {
						if ([]rune(part))[j] > 127 {
							// we replace non-ASCII char with a temporary placeholder
							// we need this to make sure size of hostname is not
							// broken by replacing non-ASCII by nothing
							newpart += "x"
						} else {
							newpart += string(([]rune(part))[j])
						}
					}
					//			// we test again with ASCII char only
					if !hostnamePartPatternRe.MatchString(newpart) {
						var validParts = hostparts[0:i]
						var notHost = hostparts[i+1:]
						var bit = hostnamePartStart.FindAllString(part, -1)
						if len(bit) > 0 {
							validParts = append(validParts, bit[1])
							notHost = append([]string{bit[2]}, notHost...)
						}
						if len(notHost) > 0 {
							rest = strings.Join(notHost[:], ".") + rest
						}
						retUrl.Hostname = strings.Join(validParts[:], ".")
						break
					}
				}
			}
		}

		if utf8.RuneCountInString(retUrl.Hostname) > hostnameMaxLen {
			retUrl.Hostname = ""
		}

		//// strip [ and ] from the hostname
		//// the host field still retains them, though
		if ipv6Hostname {
			retUrl.Hostname = string(([]rune(retUrl.Hostname))[1 : utf8.RuneCountInString(retUrl.Hostname)-2])
		}
	}

	// chop off from the tail first.
	var hash = strings.Index(rest, "#")
	if hash != -1 {
		// got a fragment string.
		retUrl.Hash = string(([]rune(rest))[hash:])
		rest = string(([]rune(rest))[0:hash])
	}

	var qm = strings.Index(rest, "?")
	if qm != -1 {
		retUrl.Search = string(([]rune(rest))[qm:])
		rest = string(([]rune(rest))[0:qm])
	}
	if utf8.RuneCountInString(rest) > 0 {
		retUrl.Pathname = rest
	}
	if slashedProtocol[lowerProto] &&
		utf8.RuneCountInString(retUrl.Hostname) > 0 &&
		utf8.RuneCountInString(retUrl.Pathname) == 0 {
		retUrl.Pathname = ""
	}

	return retUrl
}

func (md *MdUrl) ParseHost(host string, url *Url) {
	var port = portPatternRe.FindAllString(host, -1)

	if len(port) > 0 {
		port_ := port[0]

		if port_ == ":" {
			url.Port = string(([]rune(port_))[1:])
		}
		hostLen := utf8.RuneCountInString(host)
		portLen := utf8.RuneCountInString(port_)

		host = string(([]rune(host))[hostLen:portLen])
	}

	if utf8.RuneCountInString(host) > 0 {
		url.Hostname = host
	}
}

func (md *MdUrl) Format(url Url) string {
	var result = ""

	result += url.Protocol
	if url.Slashes {
		result += "//"
	}

	if utf8.RuneCountInString(url.Auth) > 0 {
		result += url.Auth + "@"
	}

	if utf8.RuneCountInString(url.Hostname) > 0 &&
		strings.Index(url.Hostname, ":") != -1 {
		// ipv6 address
		result += "[" + url.Hostname + "]'"
	} else {
		result += url.Hostname
	}

	if utf8.RuneCountInString(url.Port) > 0 {
		result += ":" + url.Port
	}

	result += url.Pathname
	result += url.Search
	result += url.Hash

	return result
}

func GetEncodeCache(exclude string) []string {
	cache := encodeCache
	if len(cache) > 0 {
		return cache
	}

	for i := 0; i < 128; i++ {
		ch := string(rune(i))

		if alphaCharsRe.MatchString(ch) {
			// always allow unencoded alphanumeric characters
			cache = append(cache, ch)
		} else {
			num := strings.ToUpper(strconv.FormatInt(int64(i), 16))
			num = string(([]rune(num))[0 : utf8.RuneCountInString(num)-2])
			cache = append(cache, "%"+"0"+num)
		}
	}

	for i := 0; i < utf8.RuneCountInString(exclude); i++ {
		//cache[exclude.charCodeAt(i)] = exclude[i];
		cache[([]rune(exclude))[i]] = string(([]rune(exclude))[i])
	}

	return cache
}

func EncodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	r = strings.Replace(r, "+", "%20", -1)
	return r
}

// Encode unsafe characters with percent-encoding, skipping already
// encoded sequences.
//
//  - string       - string to encode
//  - exclude      - list of characters to ignore (in addition to a-zA-Z0-9)
//  - keepEscaped  - don't encode '%' in a correct escape sequence (default: true)
func (md *MdUrl) Encode(str string, exclude string, keepEscaped bool) string {
	result := ""

	if utf8.RuneCountInString(exclude) == 0 {
		exclude = defaultChars
	}

	cache := GetEncodeCache(exclude)

	l := utf8.RuneCountInString(str)
	for i := 0; i < l; i++ {
		code := ([]rune(str))[i]

		if keepEscaped && code == 0x25 /* % */ && i+2 < l {
			if twoHexRe.MatchString(string(([]rune(str))[i+1 : i+3])) {
				result += string(([]rune(str))[i : i+3])
				i += 2
				continue
			}
		}

		if code < 128 {
			result += cache[code]
			continue
		}
		//
		if code >= 0xD800 && code <= 0xDFFF {
			if code >= 0xD800 && code <= 0xDBFF && i+1 < l {

				nextCode := ([]rune(str))[i+1]
				if nextCode >= 0xDC00 && nextCode <= 0xDFFF {
					str_ := string(([]rune(str))[i]) + string(([]rune(str))[i+1])
					result += EncodeURIComponent(str_)
					i++
					continue
				}
			}
			result += "%EF%BF%BD"
			continue
		}

		str_ := string(([]rune(str))[i])
		result += EncodeURIComponent(str_)
	}
	return result
}

func (md *MdUrl) Slice(s string, start int, end int) string {
	if end <= start {
		return ""
	}

	n := utf8.RuneCountInString(s)
	if start >= n || end > n {
		return ""
	}
	return string([]rune(s)[start:end])
}
