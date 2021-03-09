package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Node interface {
	fmt.Stringer
}

type Directory struct {
	name string
	children []Node
}

type File struct {
	name string
	size int64
}

func (directory Directory) String() string {
	return directory.name
}

func (file File) String() string {
	if file.size == 0 {
		return file.name + " (empty)"
	} else {
		return file.name + " (" + strconv.FormatInt(file.size, 10) + "b)"
	}
}

func readDirectory(path string, nodes []Node) ([]Node, error) {
	files, err := ioutil.ReadDir(path)

	for _, file := range files {
		var node Node

		if file.IsDir() {
			children, _ := readDirectory(filepath.Join(path, file.Name()), []Node{})
			node = Directory{
				name: file.Name(),
				children: children,
			}
		} else {
			node = File{
				name: file.Name(),
				size: file.Size(),
			}
		}

		nodes = append(nodes, node)
	}

	return nodes, err
}

func printDirectory(out io.Writer, nodes []Node, prefixes []string, printFiles bool) {

	fmt.Fprintf(out, "%s", strings.Join(prefixes, ""))
	if len(nodes) == 1 {
		fmt.Fprintf(out, "└───")
	} else {
		fmt.Fprintf(out, "├───")
	}

	for _, node := range nodes {
		if openedNode, ok := node.(Directory); ok {
			fmt.Fprintf(out, "%s\n", node)
			prefixes = append(prefixes, "│\t")
			printDirectory(out, openedNode.children, prefixes, printFiles)
			prefixes = prefixes[1:]
		}
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	nodes, err := readDirectory(path, []Node{})
	printDirectory(out, nodes, []string{}, printFiles)
	return err
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
