package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Node interface {
	fmt.Stringer
}

type Folder struct {
	name	string
	child	[]Node
}

type File struct {
	name 	string
	size	int64
}

func (file File) String() string {
	if file.size == 0 {
		return file.name + " (empty)"
	}
	return file.name + " (" + strconv.FormatInt(file.size, 10) + "b)"
}

func (folder Folder) String() string {
	return  folder.name
}

func readFolder(pathToFolder string, nodes []Node, printFiles bool) (error, []Node) {
	file, err := os.Open(pathToFolder)
	files, err := file.Readdir(0)
	file.Close()

	for _, file := range files {
		if !(file.IsDir() || printFiles) {
			continue
		}

		var newNode Node
		if file.IsDir() {
			 _, children := readFolder(filepath.Join(pathToFolder, file.Name()), []Node{}, printFiles)
			newNode = Folder{ name:  file.Name(), child: children,
			}
		} else {
			newNode = File{ file.Name(), file.Size() }
		}

		nodes = append(nodes, newNode)
	}

	return err, nodes
}

func printFolder(out io.Writer, nodes []Node, prefixes []string) {
	if len(nodes) == 0 {
		return
	}

	fmt.Fprintf(out, "%s", strings.Join(prefixes, ""))

	node := nodes[0]

	if len(nodes) == 1 {
		fmt.Fprintf(out, "%s%s\n", "└───", node)
		if folder, ok := node.(Folder); ok {
			printFolder(out, folder.child, append(prefixes, "\t"))
		}
		return
	}

	fmt.Fprintf(out, "%s%s\n", "├───", node)
	if folder, ok := node.(Folder); ok {
		printFolder(out, folder.child, append(prefixes, "│\t"))
	}

	printFolder(out, nodes[1:], prefixes)
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err, nodes := readFolder(path, []Node{}, printFiles)
	printFolder(out, nodes, []string{})

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