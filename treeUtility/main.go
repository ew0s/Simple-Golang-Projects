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

func readDirectory(path string, nodes []Node, printFiles bool) ([]Node, error) {
	files, err := ioutil.ReadDir(path)

	for _, file := range files {
		if !(file.IsDir() || printFiles) {
			continue
		}
		var node Node

		if file.IsDir() {
			children, _ := readDirectory(filepath.Join(path, file.Name()), []Node{}, printFiles)
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

func printDirectory(out io.Writer, nodes []Node, prefixes []string) {

	if len(nodes) == 0 {
		return
	}

	fmt.Fprintf(out, "%s", strings.Join(prefixes, ""))
	node := nodes[0]

	if len(nodes) == 1 {
		fmt.Fprintf(out, "%s%s\n", "└───", node)
		if directory, ok := node.(Directory); ok {
			printDirectory(out, directory.children, append(prefixes, "\t"))
		}
		return
	}

	fmt.Fprintf(out, "%s%s\n", "├───", node)
	if directory, ok := node.(Directory); ok {
		printDirectory(out, directory.children, append(prefixes, "│\t"))
	}

	printDirectory(out, nodes[1:], prefixes)
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	nodes, err := readDirectory(path, []Node{}, printFiles)
	printDirectory(out, nodes, []string{})
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
