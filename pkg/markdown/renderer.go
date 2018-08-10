package markdown

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/golang/glog"
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
	glog.V(4).Infof("%T (%v): %+v\n", node, entering, node)
	switch node := node.(type) {
	case *ast.Document:
		r.document(w, node, entering)
	case *ast.Softbreak:
		r.unsupported(node)
	case *ast.Hardbreak:
		r.hardBreak(w, node)
	case *ast.Emph:
		r.unsupported(node)
	case *ast.Strong:
		r.unsupported(node)
	case *ast.Del:
		r.unsupported(node)
	case *ast.BlockQuote:
		r.unsupported(node)
	case *ast.Aside:
		r.unsupported(node)
	case *ast.Link:
		r.link(w, node, entering)
	case *ast.CrossReference:
		r.unsupported(node)
	case *ast.Citation:
		r.unsupported(node)
	case *ast.Image:
		r.unsupported(node)
	case *ast.Code:
		r.unsupported(node)
	case *ast.CodeBlock:
		r.codeBlock(w, node, entering)
	case *ast.Caption:
		r.unsupported(node)
	case *ast.CaptionFigure:
		r.unsupported(node)
	case *ast.Paragraph:
		r.paragraph(w, node, entering)
	case *ast.HTMLSpan:
		r.unsupported(node)
	case *ast.HTMLBlock:
		r.unsupported(node)
	case *ast.Heading:
		r.heading(w, node, entering)
	case *ast.Text:
		r.text(w, node)
	case *ast.HorizontalRule:
		r.unsupported(node)
	case *ast.List:
		r.list(w, node, entering)
	case *ast.ListItem:
		r.listItem(w, node, entering)
	case *ast.Table:
		r.unsupported(node)
	case *ast.TableCell:
		r.unsupported(node)
	case *ast.TableHead:
		r.unsupported(node)
	case *ast.TableBody:
		r.unsupported(node)
	case *ast.TableRow:
		r.unsupported(node)
	case *ast.Math:
		r.unsupported(node)
	case *ast.MathBlock:
		r.unsupported(node)
	case *ast.DocumentMatter:
		r.unsupported(node)
	default:
		panic(fmt.Sprintf("Unknown node %T: %+v", node, node))
	}
	return ast.GoToNext
}

func (r *Renderer) unsupported(node ast.Node) {
	panic(fmt.Sprintf("Unsupported node %T: %+v", node, node))
}

func (r *Renderer) RenderHeader(w io.Writer, node ast.Node) {
	return
}

func (r *Renderer) RenderFooter(w io.Writer, node ast.Node) {
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

func (r *Renderer) link(w io.Writer, link *ast.Link, entering bool) {
	if entering {
		w.Write([]byte("["))
	} else {
		fmt.Fprintf(w, "](%s)", link.Destination)
	}
}

func (r *Renderer) codeBlock(w io.Writer, codeBlock *ast.CodeBlock, entering bool) {
	if entering {
		r.openSection(w)
		fmt.Fprintf(w, "```%s\n", codeBlock.Info)
		w.Write(codeBlock.Literal)
		fmt.Fprintf(w, "```\n")
	}
}

func NewRenderer() markdown.Renderer {
	return &Renderer{
		lastLine: make([]byte, 0),
	}
}
