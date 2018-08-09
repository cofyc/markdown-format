package markdown

import (
	"io"

	"github.com/gomarkdown/markdown/ast"
)

// Renderer implements Renderer interface for markdown output.
type Renderer struct {
}

func (r *Renderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	if entering == false {
		return ast.GoToNext
	}
	switch node := node.(type) {
	case *ast.Text:
	}
	return ast.GoToNext
}
