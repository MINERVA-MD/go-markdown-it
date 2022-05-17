package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-markdown-it/pkg"
	"path/filepath"
	"runtime/debug"
	"strings"
	"testing"
)

func FilterFixtures(fixtures []Fixture, filter string) []Fixture {
	if filter == "" {
		return fixtures
	}

	var filteredFixtures []Fixture

	for _, fixture := range fixtures {
		if strings.Contains(fixture.Header, filter) {
			filteredFixtures = append(filteredFixtures, fixture)
		}
	}

	return filteredFixtures
}

func TestGoMarkdownIt(t *testing.T) {
	result, err := GetResults(filepath.Join("fixtures", "markdown-it", "tables.txt"))

	if err != nil {
		fmt.Println(err)
	} else {
		// Process Results
		specs := result.Fixtures

		fix := "GFM 4.10 Tables (extension), Example 203"
		specs = FilterFixtures(specs, fix)

		var md = &pkg.MarkdownIt{}
		_ = md.MarkdownIt("default", pkg.Options{Html: true, LangPrefix: "", Typography: true, Linkify: true})

		for _, spec := range specs {
			specTitle := "Spec: " + spec.Header + "| " + spec.Type
			t.Run(specTitle, func(t *testing.T) {
				if fix == "" {
					defer func() {
						if err := recover(); err != nil {
							fmt.Println("stacktrace from test: " + specTitle + "\n" + string(debug.Stack()))
							t.Fail()
						}
					}()
				}
				assert.Equal(t, Normalize(spec.Second.Text), md.Render(spec.First.Text, &pkg.Env{}))
			})
		}
	}
}
