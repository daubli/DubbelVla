package main

import (
	"os"
	"fmt"
	"path/filepath"
	"path"
	"runtime"
	"log"
	"math"
	"io"
)

var FilenameMap map[string][]string
const filechunk = 8192

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

func md5(filePath string)[]byte{
	file, error := os.Open(filePath)
	if error != nil {
		return nil
	}
	defer file.Close()

	info, _ := file.Stat()
	size := info.Size()

	blocks := uint64(math.Ceil(float64(size) / float64(filechunk)))

	hashsum := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(size-int64(i*filechunk))))
		buf := make([] byte, blocksize)
		file.Read(buf)
		io.WriteString(hashsum, string(buf))   // append into the hash
	}

	return hashsum.Sum(nil)
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
