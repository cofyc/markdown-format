package markdown

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
)

// Renderer implements Renderer interface for markdown output.
type Renderer struct {
	documentBegin bool
	lastLine      []byte
}

type WriterFunc func(p []byte) (n int, err error)

func (f WriterFunc) Write(p []byte) (n int, err error) {
	return f(p)
}

func (f *Renderer) writerWrapper(w io.Writer) io.Writer {
	return WriterFunc(func(p []byte) (n int, err error) {
		if len(p) > 0 {
			f.lastLine = append(f.lastLine, p...)
			seps := bytes.Split(f.lastLine, []byte{'\n'})
			if len(seps) == 1 {
				f.lastLine = seps[0]
			} else {
				if len(seps[len(seps)-1]) == 0 {
					f.lastLine = append(seps[len(seps)-2], '\n')
				}
			}
		}
		return w.Write(p)
	})
}

func (r *Renderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	w = r.writerWrapper(w)
	fmt.Printf("%T (%v): %v\n", node, entering, node)
	switch node := node.(type) {
	case *ast.Document:
		r.document(w, node, entering)
	case *ast.Softbreak:
		// TODO
	case *ast.Hardbreak:
		r.hardBreak(w, node)
	case *ast.Emph:
		// TODO
	case *ast.Strong:
		// TODO
	case *ast.Del:
		// TODO
	case *ast.BlockQuote:
		// TODO
	case *ast.Aside:
		// TODO
	case *ast.Link:
		// TODO
	case *ast.CrossReference:
		// TODO
	case *ast.Citation:
		// TODO
	case *ast.Image:
		// TODO
	case *ast.Code:
		// TODO
	case *ast.CodeBlock:
		// TODO
	case *ast.Caption:
		// TODO
	case *ast.CaptionFigure:
		// TODO
	case *ast.Paragraph:
		r.paragraph(w, node, entering)
	case *ast.HTMLSpan:
		// TODO
	case *ast.HTMLBlock:
		// TODO
	case *ast.Heading:
		r.heading(w, node, entering)
	case *ast.Text:
		r.text(w, node)
	case *ast.HorizontalRule:
		// TODO
	case *ast.List:
		r.list(w, node, entering)
	case *ast.ListItem:
		r.listItem(w, node, entering)
	case *ast.Table:
		// TODO
	case *ast.TableCell:
		// TODO
	case *ast.TableHead:
		// TODO
	case *ast.TableBody:
		// TODO
	case *ast.TableRow:
		// TODO
	case *ast.Math:
		// TODO
	case *ast.MathBlock:
		// TODO
	case *ast.DocumentMatter:
		// TODO
	default:
		panic(fmt.Sprintf("Unknown node %T", node))
	}
	return ast.GoToNext
}

func (r *Renderer) RenderHeader(w io.Writer, ast ast.Node) {
	return
}

func (r *Renderer) RenderFooter(w io.Writer, ast ast.Node) {
	return
}

func (r *Renderer) document(w io.Writer, doc *ast.Document, entering bool) {
	r.documentBegin = entering
}

func (r *Renderer) heading(w io.Writer, heading *ast.Heading, entering bool) {
	if entering {
		r.openSection(w)
		fmt.Fprintf(w, "%s ", strings.Repeat("#", heading.Level))
	} else {
		r.closeSection(w)
	}
}

func (r *Renderer) openSection(w io.Writer) {
	// Add newline to open section if it is not at begin of document, or last line is not a newline.
	if !r.documentBegin && bytes.Compare(r.lastLine, []byte("\n")) != 0 {
		w.Write([]byte("\n"))
	}
	r.documentBegin = false
}

func (r *Renderer) closeSection(w io.Writer) {
	// Add newline to end section.
	if bytes.Compare(r.lastLine, []byte("\n")) != 0 {
		w.Write([]byte("\n"))
	}
}

func (r *Renderer) list(w io.Writer, list *ast.List, entering bool) {
	if entering {
		r.openSection(w)
	} else {
		// r.closeSection(w)
	}
}

func (r *Renderer) listItem(w io.Writer, listItem *ast.ListItem, entering bool) {
	if entering {
		fmt.Fprintf(w, "%c ", listItem.BulletChar)
	} else {
		r.closeSection(w)
	}
}

func (r *Renderer) text(w io.Writer, text *ast.Text) {
	w.Write(text.Literal)
}

func (r *Renderer) paragraph(w io.Writer, paragraph *ast.Paragraph, entering bool) {
	if _, ok := paragraph.Parent.(*ast.ListItem); ok {
		return
	}
	if entering {
		r.openSection(w)
	} else {
		r.closeSection(w)
	}
}

func (r *Renderer) hardBreak(w io.Writer, node *ast.Hardbreak) {
	w.Write([]byte("\n"))
}

func NewRenderer() markdown.Renderer {
	return &Renderer{
		lastLine: make([]byte, 0),
	}
}
