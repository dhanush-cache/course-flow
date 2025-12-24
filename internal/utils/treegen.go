package utils

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
)

var (
	dirColor     = color.New(color.FgBlue, color.Bold)
	fileColor    = color.New(color.FgGreen)
	connectorCol = color.New(color.FgHiBlack)
)

type Node struct {
	Name     string
	Children map[string]*Node
	IsFile   bool
}

func NewNode(name string, isFile bool) *Node {
	return &Node{
		Name:     name,
		Children: make(map[string]*Node),
		IsFile:   isFile,
	}
}

func BuildTree(paths []string) []string {
	root := NewNode("", false)

	for _, path := range paths {
		parts := strings.Split(path, string(filepath.Separator))
		current := root

		for i, part := range parts {
			isFile := i == len(parts)-1
			if _, ok := current.Children[part]; !ok {
				current.Children[part] = NewNode(part, isFile)
			}
			current = current.Children[part]
		}
	}

	var output []string
	printTree(root, "", &output)
	return output
}

func printTree(node *Node, prefix string, output *[]string) {
	keys := make([]string, 0, len(node.Children))
	for k := range node.Children {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i, key := range keys {
		child := node.Children[key]
		last := i == len(keys)-1

		var connector string
		if last {
			connector = "└── "
		} else {
			connector = "├── "
		}

		coloredConnector := connectorCol.Sprint(connector)
		coloredPrefix := connectorCol.Sprint(prefix)

		var coloredName string
		if child.IsFile {
			coloredName = fileColor.Sprint(child.Name)
		} else {
			coloredName = dirColor.Sprint(child.Name)
		}

		line := coloredPrefix + coloredConnector + coloredName
		*output = append(*output, line)

		var newPrefix string
		if last {
			newPrefix = prefix + "    "
		} else {
			newPrefix = prefix + "│   "
		}

		printTree(child, newPrefix, output)
	}
}
