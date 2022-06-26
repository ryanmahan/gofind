package routines

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

func handleDir(path string, channel chan<- map[string]string) {
	entries, err := os.ReadDir(path)
	check(err)
	for _, entry := range entries {
		entryInfo, err := entry.Info()
		check(err)
		if entryInfo.IsDir() {
			handleDir(path+"/"+entryInfo.Name(), channel)
		} else {
			indexFile(path+"/"+entryInfo.Name(), channel)
		}
	}
}

func indexFile(path string, channel chan<- map[string]string) {
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

	index := make(map[string]string)
	reg, err := regexp.Compile("(^[^A-Za-z0-9]+)|([^A-Za-z0-9]+$)")

	for scanner.Scan() {
		word := reg.ReplaceAllString(strings.ToLower(scanner.Text()), "")
		if len(word) > 0 {
			index[word] = path
		}
	}
	channel <- index
}

func main() {

	args := os.Args[1:]
	path := args[0]
	file, err := os.Open(args[0])
	check(err)
	fileInfo, err := file.Stat()
	check(err)
	index := make(map[string][]string)
	size := 1

	if fileInfo.IsDir() {
		dir, err := os.ReadDir(path)
		check(err)
		size = len(dir)
	}

	channel := make(chan map[string]string, size)
	handleDir(path, channel)

	for i := 0; i < size; i++ {
		fileIndex := <-channel
		for key, path := range fileIndex {
			indexed, present := index[key]
			if present {
				indexed = append(indexed, path)
				index[key] = indexed
			} else {
				index[key] = []string{path}
			}
		}
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
