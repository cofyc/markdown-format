package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

func toSectionLink(name string) string {
	name = strings.ToLower(name)
	name = strings.Replace(name, " ", "-", -1)
	return name
}

type TOCGenerator struct {
	tocNode     *ast.Container
	parentLevel int
	parentNode  ast.Node
	lastNode    ast.Node
	nodeStack   []ast.Node
}

func (g *TOCGenerator) Visit(node ast.Node, entering bool) ast.WalkStatus {
	if entering == false {
		return ast.GoToNext
	}
	switch n := node.(type) {
	case *ast.Heading:
		if n.Level == 1 {
			break
		}
		if g.parentNode == nil {
			g.parentNode = g.tocNode
		}
		if n.Level > g.parentLevel {
			// down
			if n.Level-g.parentLevel > 1 {
				glog.Fatal("TODO")
			}
			g.nodeStack = append(g.nodeStack, g.parentNode)
			g.parentNode = g.lastNode
			g.parentLevel = n.Level
		} else if n.Level < g.parentLevel {
			// up
			for i := 0; i < g.parentLevel-n.Level; i++ {
				g.parentNode, g.nodeStack = g.nodeStack[len(g.nodeStack)-1], g.nodeStack[:len(g.nodeStack)-1]
			}
			g.parentLevel = n.Level
		}
		listNode := &ast.ListItem{}
		if n.Level == 2 {
			listNode.BulletChar = '-'
		} else {
			listNode.BulletChar = '*'
		}
		if len(n.Children) <= 0 {
			glog.Fatal("TODO heading node has no children")
		}
		if textChild, ok := n.Children[0].(*ast.Text); !ok {
			glog.Errorf("heading node contains a non-text node")
		} else {
			listNode.Literal = []byte(fmt.Sprintf("[%s](#%s)", textChild.Literal, n.HeadingID))
		}
		parentContainer := g.parentNode.AsContainer()
		parentContainer.Children = append(parentContainer.Children, listNode)
		g.lastNode = listNode
	}
	return ast.GoToNext
}

func toc(node ast.Node) ast.Node {
	tocGenerator := &TOCGenerator{
		tocNode:     &ast.Container{},
		parentLevel: 2,
	}
	ast.Walk(node, tocGenerator)
	return tocGenerator.tocNode
}

func printTOC(node ast.Node, buf *bytes.Buffer, indent int) {
	if len(node.GetChildren()) == 0 {
		return
	}
	if indent > 5 {
		glog.Fatal("Someting goes wrong!")
	}
	for _, child := range node.GetChildren() {
		listItem := child.(*ast.ListItem)
		fmt.Fprintf(buf, "%s%c %s\n", strings.Repeat("  ", indent), listItem.BulletChar, string(listItem.Literal))
		printTOC(listItem, buf, indent+1)
	}
}

func generate(data []byte) ([]byte, error) {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	node := markdown.Parse(data, parser)
	tocNode := toc(node)
	var buf bytes.Buffer
	printTOC(tocNode, &buf, 0)
	return buf.Bytes(), nil
}

func process(f string) {
	fdata, err := ioutil.ReadFile(f)
	if err != nil {
		glog.Fatal(err)
	}
	ftocdata, err := generate(fdata)
	if err != nil {
		glog.Fatal(err)
	}
	fmt.Fprintf(os.Stdout, "%s", ftocdata)
}

func main() {
	filenames := os.Args[1:]
	if len(filenames) == 1 {
		process(filenames[0])
	} else {
		for _, f := range os.Args[1:] {
			fmt.Fprintf(os.Stdout, "### %s\n", f)
			process(f)
		}
	}
}
