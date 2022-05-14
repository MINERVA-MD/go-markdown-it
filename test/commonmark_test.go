package test

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-markdown-it/pkg"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strconv"
	"testing"
)

type Spec struct {
	Markdown   string
	Html       string
	Example    int
	StartLine  int
	EndLine    int
	Section    string
	ShouldFail bool
	Marked     string
	MarkdownIt string
}

func GetCurrPath() (string, error) {
	currDir, err := os.Getwd()

	if err != nil {
		return "", err
	}

	return currDir, nil
}

func GetResults(relativePath string) (Result, error) {
	base, err := GetCurrPath()

	if err == nil {
		path := filepath.Join(base, relativePath)
		dat, err := os.ReadFile(path)

		if err != nil {
			return Result{}, err
		} else {
			results := Parse(string(dat), Options{Sep: []string{"."}})
			return results, nil
		}
	}

	return Result{}, err
}

func ReadFileContents(relativePath string) (string, error) {
	base, err := GetCurrPath()

	if err == nil {
		path := filepath.Join(base, relativePath)
		dat, err := os.ReadFile(path)

		if err != nil {
			return "", err
		} else {
			return string(dat), nil
		}
	}
	return "", err
}

func Normalize(text string) string {
	BlockquoteRe := regexp.MustCompile("<blockquote>\n</blockquote>")
	return BlockquoteRe.ReplaceAllString(text, "<blockquote></blockquote>")
}

func GetSpec(specs []Spec, example int) []Spec {
	if example == -1 {
		return specs
	}

	var filteredSpecs []Spec

	for _, spec := range specs {
		if spec.Example == example {
			filteredSpecs = append(filteredSpecs, spec)
		}
	}

	return filteredSpecs
}

func TestCommonMark(t *testing.T) {
	//results, err := GetResults(filepath.Join("fixtures", "commonmark", "good.txt"))

	specsJson, err := ReadFileContents(filepath.Join("fixtures", "commonmark.0.30.json"))

	if err != nil {
		fmt.Println(err)
	} else {
		// Process Results
		var md = &pkg.MarkdownIt{}
		_ = md.MarkdownIt("commonmark", pkg.Options{Html: true, XhtmlOut: true, LangPrefix: "language-"})

		var specs []Spec
		_ = json.Unmarshal([]byte(specsJson), &specs)

		// code_inline
		num := -1
		specs = GetSpec(specs, num)

		total := 0
		passed := 0
		failed := 0

		for _, spec := range specs {
			specTitle := "Spec: " + strconv.Itoa(spec.Example) + "| " + spec.Section
			t.Run(specTitle, func(t *testing.T) {
				if num == -1 {
					defer func() {
						if err := recover(); err != nil {
							passed--
							failed++
							fmt.Println("Recovered in f", err)
							fmt.Println("stacktrace from test: " + specTitle + "\n" + string(debug.Stack()))
							t.Fail()
						}
					}()
				}

				ok := assert.Equal(t, Normalize(spec.Html), md.Render(spec.Markdown, &pkg.Env{}))

				if ok {
					passed++
				} else {
					failed++
				}
			})
			total++
		}
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println(passed, "/", total, "=", (float32(passed)/float32(total))*float32(100), "%")
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
	}
}

//func TestSpec472(t *testing.T) {
//	var md = &pkg.MarkdownIt{}
//	_ = md.MarkdownIt("commonmark", pkg.Options{Html: true, XhtmlOut: true})
//
//	assert.Equal(t, "<hr />\n", md.Render("*\t*\t*\t\n", pkg.Env{}))
//}
//
//func TestSpec489(t *testing.T) {
//	var md = &pkg.MarkdownIt{}
//	_ = md.MarkdownIt("commonmark", pkg.Options{Html: true, XhtmlOut: true})
//
//	assert.Equal(t, "<p>!&quot;#$%&amp;'()*+,-./:;&lt;=&gt;?@[\\]^_`{|}~</p>", md.Render("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~\n", pkg.Env{}))
//}
//
//func TestSpec499(t *testing.T) {
//	var md = &pkg.MarkdownIt{}
//	_ = md.MarkdownIt("commonmark", pkg.Options{Html: true, XhtmlOut: true})
//
//	assert.Equal(t, "<p>\\\t\\A\\a\\ \\3\\φ\\«</p>\n", md.Render("\\\t\\A\\a\\ \\3\\φ\\«\n", pkg.Env{}))
//}
//
//func TestSpec509(t *testing.T) {
//	var md = &pkg.MarkdownIt{}
//	_ = md.MarkdownIt("commonmark", pkg.Options{Html: true, XhtmlOut: true})
//
//	assert.Equal(t, "<p>*not emphasized*\n&lt;br/&gt; not a tag\n[not a link](/foo)\n`not code`\n1. not a list\n* not a list\n# not a heading\n[foo]: /url &quot;not a reference&quot;\n&amp;ouml; not a character entity</p>\n", md.Render("\\*not emphasized*\n\\<br/> not a tag\n\\[not a link](/foo)\n\\`not code`\n1\\. not a list\n\\* not a list\n\\# not a heading\n\\[foo]: /url \"not a reference\"\n\\&ouml; not a character entity\n", pkg.Env{}))
//}
//
//func TestSpec534(t *testing.T) {
//	var md = &pkg.MarkdownIt{}
//	_ = md.MarkdownIt("commonmark", pkg.Options{Html: true, XhtmlOut: true})
//
//	assert.Equal(t, "<p>\\<em>emphasis</em></p>\n", md.Render("\\\\*emphasis*\n", pkg.Env{}))
//}
//
//func TestSpec555(t *testing.T) {
//	var md = &pkg.MarkdownIt{}
//	_ = md.MarkdownIt("commonmark", pkg.Options{Html: true, XhtmlOut: true})
//
//	assert.Equal(t, "<p><code>\\[\\`</code></p>\n", md.Render("`` \\[\\` ``\n", pkg.Env{}))
//}
