package main

import (
	"log"
	"os"
	"sort"
	"strings"

	"github.com/ryanmahan/gofind/common"
	"github.com/ryanmahan/gofind/readers"
)

func fileRouter(path string) readers.Tokenizer {
	log.Println("Processing " + path)
	if strings.Contains(path, ".pdf") {
		return readers.NewPDFTokenizer(path)
	} else {
		return readers.NewBaseTokenizer(path)
	}
}

func handleDir(path string) map[string][]string {
	entries, err := os.ReadDir(path)
	common.Check(err)
	index := make(map[string][]string)
	tokenChannel := make(chan readers.FileTokens)
	channelEntries := 0
	for _, entry := range entries {
		entryInfo, err := entry.Info()
		common.Check(err)
		filepath := path + "/" + entryInfo.Name()
		tokenizer := fileRouter(filepath)
		go tokenizer.Index(tokenChannel)
		channelEntries++
	}

	for i := 0; i < channelEntries; i++ {
		tokenized := <-tokenChannel
		for _, word := range tokenized.TokenList {
			if index[word] != nil {
				index[word] = append(index[word], tokenized.Path)
			} else {
				index[word] = []string{tokenized.Path}
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
