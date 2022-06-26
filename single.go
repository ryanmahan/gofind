package main

import (
	"bufio"
	"os"
	"regexp"
	"sort"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func handleDir(path string) map[string][]string {
	entries, err := os.ReadDir(path)
	check(err)
	index := make(map[string][]string)
	for _, entry := range entries {
		entryInfo, err := entry.Info()
		check(err)
		curr := make(map[string][]string)
		if entryInfo.IsDir() {
			curr = handleDir(path + "/" + entryInfo.Name())
		} else {
			curr = indexFile(path + "/" + entryInfo.Name())
		}

		for key, entry := range curr {
			if index[key] != nil {
				index[key] = append(index[key], entry...)
			} else {
				index[key] = entry
			}
		}
	}
	return index
}

func indexFile(path string) map[string][]string {
	file, err := os.Open(path)
	check(err)
	scanner := bufio.NewScanner(file)

	scanner.Split(func(data []byte, eof bool) (advance int, token []byte, err error) {
		if eof {
			return len(data), data, bufio.ErrFinalToken
		}
		for index, token := range data {
			if index > 0 && (token == '\n' || token == ' ') {
				return index + 1, data[0:index], nil
			}
		}
		return 0, nil, nil
	})

	index := make(map[string][]string)
	reg, err := regexp.Compile("(^[^A-Za-z0-9]+)|([^A-Za-z0-9]+$)")
	check(err)

	for scanner.Scan() {
		word := reg.ReplaceAllString(strings.ToLower(scanner.Text()), "")
		if len(word) > 0 {
			index[word] = []string{path}
		}
	}
	return index
}

func main() {

	args := os.Args[1:]
	path := args[0]
	file, err := os.Open(args[0])
	check(err)
	fileInfo, err := file.Stat()
	check(err)
	var index map[string][]string
	if fileInfo.IsDir() {
		index = handleDir(path)
	} else {
		index = indexFile(path)
	}

	outputFile, err := os.Create("index.txt")
	check(err)

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
