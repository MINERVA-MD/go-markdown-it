package test

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-markdown-it/pkg"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
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

func TestAgainstCommonMarkSpec(t *testing.T) {
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

func TestCommonMarkWithGoodData(t *testing.T) {
	result, err := GetResults(filepath.Join("fixtures", "commonmark", "good.txt"))

	if err != nil {
		fmt.Println(err)
	} else {
		// Process Results
		specs := result.Fixtures

		var md = &pkg.MarkdownIt{}
		_ = md.MarkdownIt("commonmark", pkg.Options{Html: true, XhtmlOut: true, LangPrefix: "language-"})

		for _, spec := range specs {
			specTitle := "Spec: " + spec.Header + "| " + spec.Type

			t.Run(specTitle, func(t *testing.T) {
				defer func() {
					if err := recover(); err != nil {
						fmt.Println("stacktrace from test: " + specTitle + "\n" + string(debug.Stack()))
						t.Fail()
					}
				}()

				assert.Equal(t, Normalize(spec.Second.Text), md.Render(spec.First.Text, &pkg.Env{}))
			})
		}
	}
}

type MDSpec struct {
	Title    string
	Markdown string
	Html     string
	Lines    int
}

func GetFullPath(relativePath string) string {
	base, err := GetCurrPath()
	if err != nil {
		panic("Unable to get current path")
	}
	return filepath.Join(base, relativePath)
}

func TestMDFilesWithCM(t *testing.T) {
	mdSpecs := []MDSpec{
		{
			Title:    "Commonmark Spec (via Markdown-It)",
			Markdown: filepath.Join("fixtures", "md", "cm", "spec.md"),
			Html:     filepath.Join("fixtures", "md", "cm", "spec.html"),
		},
	}

	for _, spec := range mdSpecs {
		// Process Results
		var md = &pkg.MarkdownIt{}
		_ = md.MarkdownIt("commonmark", pkg.Options{Html: true, XhtmlOut: true, LangPrefix: "language-"})

		specMD, err := ReadFileContents(spec.Markdown)
		specHTML, err := ReadFileContents(spec.Html)

		if err != nil {
			fmt.Println(err)
		} else {
			t.Run(spec.Title, func(t *testing.T) {
				assert.Equal(t, specHTML, md.Render(specMD, &pkg.Env{}))
			})
		}
	}
}

func GetDirFilenames(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var filenames []string

	for _, f := range files {
		name := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		if !f.IsDir() && !Contains(filenames, name) {
			filenames = append(filenames, name)
		}
	}
	return filenames
}

func TestLargeMDTables(t *testing.T) {
	mdSpecs := []MDSpec{
		{
			Title:    "Large Tables",
			Markdown: filepath.Join("fixtures", "md", "cm", "tables_spec.md"),
			Html:     filepath.Join("fixtures", "md", "cm", "tables_spec.html"),
			Lines:    7580,
		},
	}

	for _, spec := range mdSpecs {
		// Process Results
		var md = &pkg.MarkdownIt{}
		_ = md.MarkdownIt("default", pkg.Options{Html: true, Typography: true, LangPrefix: ""})

		specMD, err := ReadFileContents(spec.Markdown)
		specHTML, err := ReadFileContents(spec.Html)

		if err != nil {
			fmt.Println(err)
		} else {
			t.Run(spec.Title, func(t *testing.T) {
				actualHTML := md.Render(specMD, &pkg.Env{})
				assert.Equal(t, specHTML, actualHTML)
			})
		}
	}
}

func TestOriginalSpecs(t *testing.T) {

	base, err := GetCurrPath()
	relativePath := filepath.Join("fixtures", "md", "original")
	if err == nil {
		path := filepath.Join(base, relativePath)
		filenames := GetDirFilenames(path)

		for _, filename := range filenames {
			// Process Results
			var md = &pkg.MarkdownIt{}
			_ = md.MarkdownIt("commonmark", pkg.Options{Html: true, XhtmlOut: true, LangPrefix: "language-"})

			specMD, err := ReadFileContents(filepath.Join("fixtures", "md", "original", filename+".md"))
			specHTML, err := ReadFileContents(filepath.Join("fixtures", "md", "original", filename+".html"))

			if err != nil {
				fmt.Println(err)
			} else {
				t.Run(filename, func(t *testing.T) {
					actualHTML := md.Render(specMD, &pkg.Env{})
					assert.Equal(t, specHTML, actualHTML)
				})
			}
		}
	}
}

func BenchmarkMDParser(b *testing.B) {

	specMD, err := ReadFileContents(filepath.Join("fixtures", "commonmark", "spec.md"))

	if err != nil {
		fmt.Println(err)
	} else {
		// Process Results
		var md = &pkg.MarkdownIt{}
		_ = md.MarkdownIt("commonmark", pkg.Options{Html: true, XhtmlOut: true, LangPrefix: "language-"})

		for i := 0; i < b.N; i++ {
			md.Parse(specMD, &pkg.Env{})
		}
	}
}

func BenchmarkConvertToRunes(b *testing.B) {
	b.Skip()
	specMD, err := ReadFileContents(filepath.Join("fixtures", "commonmark", "spec.md"))

	if err != nil {
		fmt.Println(err)
	} else {
		b.Run("File contents as Bytes", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = []byte(specMD)
			}
		})

		b.Run("File contents as Runes", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = []rune(specMD)
			}
		})
	}
}

func TestSpec472(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("commonmark", pkg.Options{Html: true, XhtmlOut: true})

	assert.Equal(t, "<hr />\n", md.Render("*\t*\t*\t\n", &pkg.Env{}))
}

func TestSpec489(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("commonmark", pkg.Options{Html: true, XhtmlOut: true})

	assert.Equal(t, "<p>!&quot;#$%&amp;'()*+,-./:;&lt;=&gt;?@[]^_`{|}~</p>\n", md.Render("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~\n", &pkg.Env{}))
}

func TestSpec499(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("commonmark", pkg.Options{Html: true, XhtmlOut: true})

	assert.Equal(t, "<p>\\\t\\A\\a\\ \\3\\φ\\«</p>\n", md.Render("\\\t\\A\\a\\ \\3\\φ\\«\n", &pkg.Env{}))
}

func TestSpec509(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("commonmark", pkg.Options{Html: true, XhtmlOut: true})

	assert.Equal(t, "<p>*not emphasized*\n&lt;br/&gt; not a tag\n[not a link](/foo)\n`not code`\n1. not a list\n* not a list\n# not a heading\n[foo]: /url &quot;not a reference&quot;\n&amp;ouml; not a character entity</p>\n", md.Render("\\*not emphasized*\n\\<br/> not a tag\n\\[not a link](/foo)\n\\`not code`\n1\\. not a list\n\\* not a list\n\\# not a heading\n\\[foo]: /url \"not a reference\"\n\\&ouml; not a character entity\n", &pkg.Env{}))
}

func TestSpec534(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("commonmark", pkg.Options{Html: true, XhtmlOut: true})

	assert.Equal(t, "<p>\\<em>emphasis</em></p>\n", md.Render("\\\\*emphasis*\n", &pkg.Env{}))
}

func TestSpec555(t *testing.T) {
	var md = &pkg.MarkdownIt{}
	_ = md.MarkdownIt("commonmark", pkg.Options{Html: true, XhtmlOut: true})

	assert.Equal(t, "<p><code>\\[\\`</code></p>\n", md.Render("`` \\[\\` ``\n", &pkg.Env{}))
}
