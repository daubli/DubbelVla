package main

import (
	"os"
	"fmt"
	"path/filepath"
	"path"
	"runtime"
	"log"
)

var FilenameMap map[string][]string

func find(inputDir string) {
	visit := func (filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			fmt.Println(err)
		}
		filename := path.Base(filePath)
		FilenameMap[filename] = append(FilenameMap[filename], filePath)
		return nil
	}
	FilenameMap = make(map[string][]string)

	c := make(chan error)
	go func() { c <- filepath.Walk(inputDir, visit) }()
	err := <-c
	if err!= nil{
		log.Fatal(err)
	}
}

func printFiles(){
	filesFound := 0
	for key, value := range FilenameMap {
		if len(value) < 2 {
			continue
		}
		filesFound++
		println(key, ":")
		for _, filename := range value {
			fmt.Printf("%s\n", filename)
		}
	}
	fmt.Printf("%d files with duplicates found.", filesFound)
}

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(os.Args) != 2 {
		fmt.Printf("Usage : go run main.go directory \n ", os.Args[0])
	}else{
		inputDir := os.Args[1]
		fmt.Println("Searching for duplicate files in " + inputDir)
		find(inputDir)
		printFiles()
	}
}
