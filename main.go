package main

import (
	"fmt"
	"gitlab.com/golang-commonmark/linkify"
	"regexp"
)

func main() {
	var str = "For more information, see Chapter 3.4.5.1"
	var re = regexp.MustCompile("(?i)see (chapter \\d+(\\.\\d)*)")
	var found = re.FindStringSubmatch(str)
	fmt.Println(found)

	input := `
	Check out this link to http://google.com
You can also email support@example.com to view more.

Some more links: fsf.org http://www.gnu.org/licenses/gpl-3.0.en.html 127.0.0.1
                 localhost:80	github.com/trending?l=Go	//reddit.com/r/golang
mailto:r@golang.org some.nonexistent.host.name flibustahezeous3.onion
`
	for _, l := range linkify.Links(input) {
		fmt.Printf("Scheme: %-8s  URL: %s\n", l.Scheme, input[l.Start:l.End])
	}

	//fmt.Println(found)
	//regexp.MustCompile("[\\u201c\\u201d\\u2018\\u2019]")
	//ioutil.WriteFile("UNICODE_PUNCT_RE.txt", []byte(common.UNICODE_PUNCT_RE.String()), 0644)
	//ioutil.WriteFile("ENTITY_RE.txt", []byte(common.ENTITY_RE.String()), 0644)
	//ioutil.WriteFile("UNESCAPE_ALL_RE.txt", []byte(common.UNESCAPE_ALL_RE.String()), 0644)
	//ioutil.WriteFile("DIGITAL_ENTITY_TEST_RE.txt", []byte(common.DIGITAL_ENTITY_TEST_RE.String()), 0644)
}
