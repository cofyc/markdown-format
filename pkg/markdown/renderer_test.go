package markdown

import (
	"testing"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/stretchr/testify/assert"
)

func TestRenderer(t *testing.T) {
	cases := []struct {
		input       string
		expected    string
		expectedErr error
	}{
		{
			input: `# No TOC docs

Some description.

## 1 First section

First section content.

### 1.1 Sub section

Sub section content.

### 1.2 Sub second section

Sub second section content

## 2 Second section

Second section.

## References

- First reference here.
- Second reference here.
`,
			expected: `# No TOC docs

Some description.

## 1 First section

First section content.

### 1.1 Sub section

Sub section content.

### 1.2 Sub second section

Sub second section content

## 2 Second section

Second section.

## References

- First reference here.
- Second reference here.
`,
		},
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	renderer := NewRenderer()
	for _, v := range cases {
		inputNode := parser.Parse([]byte(v.input))
		gotBytes := markdown.Render(inputNode, renderer)
		got := string(gotBytes)
		assert.Equal(t, v.expected, got)
	}
}
