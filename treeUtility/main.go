package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type File struct {
	name	string
	size	int64
}

type Folder struct {
	name		string
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

func dirTree(out *os.File, path string, printFiles bool) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic("Unable to read path directory")
	}

	for _, file := range files {
		newPath := path+string(os.PathSeparator)+file.Name()
		_ = showTree(out, newPath, false, 0)
	}

	return nil
}

func showTree(out *os.File, path string, leaf bool, level int) error {

	fileInfo, err := os.Lstat(path)
	if err != nil {
		log.Fatal("Unable to show file info")
	}

	if fileInfo.IsDir() {
		folder := &Folder{name: fileInfo.Name()}
		folder.drawComponent(level, leaf)
	} else {
		var file = &File{
			name: fileInfo.Name(),
			size: fileInfo.Size(),
		}
		file.drawComponent(level, leaf)
	}

	files, err := ioutil.ReadDir(path)
	if err == nil {
		if len(files) != 0 {
			for index, file := range files {
				leaf = false
				if index == len(files) - 1 {
					leaf = true
				}
				if file.IsDir() {
					err := showTree(out, path+string(os.PathSeparator)+file.Name(), leaf, level + 1)
					if err != nil {
						panic(err.Error())
					}
				} else {
					err := showTree(out, path+string(os.PathSeparator)+file.Name(), leaf, level + 1)
					if err != nil {
						panic(err.Error())
					}
				}
			}
		}
	}

	return nil
}

func (file *File) drawComponent(level int, Leaf bool) {
	printTabs(level)
	if Leaf {
		fmt.Printf("%s", leaf)
	} else {
		fmt.Printf("%s", element)
	}
	fmt.Printf("%s (%db)\n", file.name, file.size)
}

func (folder *Folder) drawComponent(level int, Leaf bool) {
	printTabs(level)
	if Leaf {
		fmt.Printf("%s", leaf)
	} else {
		fmt.Printf("%s", element)
	}
	fmt.Printf("%s\n", folder.name)
}

const (
	leaf =		"└───"
	element =	"├───"
)

func printTabs(level int) {
	for i := 0; i < level; i++ {
		fmt.Print("│\t")
	}
}