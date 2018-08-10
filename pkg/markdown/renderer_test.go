package markdown

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/stretchr/testify/assert"
)

func TestRenderer(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{}

	inputFiles, err := ioutil.ReadDir("testdata/input")
	if err != nil {
		t.Fatal(err)
	}
	for _, inputFile := range inputFiles {
		inputFilename := filepath.Join("testdata/input", inputFile.Name())
		inputData, err := ioutil.ReadFile(inputFilename)
		if err != nil {
			t.Fatal(err)
		}
		expectedFilename := filepath.Join("testdata/expected", inputFile.Name())
		expectedData, err := ioutil.ReadFile(expectedFilename)
		if err != nil {
			t.Fatal(err)
		}
		cases = append(cases, struct {
			name     string
			input    string
			expected string
		}{
			name:     fmt.Sprintf("format-%s", inputFilename),
			input:    string(inputData),
			expected: string(expectedData),
		})
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			extensions := parser.CommonExtensions | parser.AutoHeadingIDs
			parser := parser.NewWithExtensions(extensions)
			renderer := NewRenderer()
			inputNode := parser.Parse([]byte(v.input))
			gotBytes := markdown.Render(inputNode, renderer)
			got := string(gotBytes)
			assert.Equal(t, v.expected, got)
		})
	}
}
