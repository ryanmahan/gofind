package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/ryanmahan/gofind/common"
	"github.com/ryanmahan/gofind/readers"
)

func handleDir(path string) map[string][]string {
	entries, err := os.ReadDir(path)
	common.Check(err)
	index := make(map[string][]string)
	for _, entry := range entries {
		entryInfo, err := entry.Info()
		common.Check(err)
		tokenizer := readers.BaseTokenizer{Path: path + "/" + entryInfo.Name()}
		curr := tokenizer.Index()
		for _, word := range curr {
			fmt.Println(word)
			if index[word] != nil {
				index[word] = append(index[word], tokenizer.Path)
			} else {
				index[word] = []string{tokenizer.Path}
			}
		}
	}
	return index
}

func main() {

	args := os.Args[1:]
	path := args[0]
	file, err := os.Open(args[0])
	common.Check(err)
	fileInfo, err := file.Stat()
	common.Check(err)
	index := make(map[string][]string)
	if fileInfo.IsDir() {
		index = handleDir(path)
	} else {
		var tokenizer readers.Tokenizer
		if strings.Contains(path, ".pdf") {
			tokenizer = readers.NewPDFTokenizer(path)
		} else {
			tokenizer = readers.NewBaseTokenizer(path)
		}
		
		curr := tokenizer.Index()
		for _, word := range curr {
			fmt.Println(word)
			if index[word] != nil {
				index[word] = append(index[word], path)
			} else {
				index[word] = []string{path}
			}
		}
	}

	outputFile, err := os.Create("index.txt")
	common.Check(err)

	keys := make([]string, len(index))

	for k := range index {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		if len(key) > 0 {
			outputFile.WriteString(key + "\t " + strings.Join(index[key], ", ") + "\n")
		}
	}
}
