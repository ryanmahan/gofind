package readers

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/ryanmahan/gofind/common"
)

type Tokenizer interface {
	Index(tokens chan FileTokens)
}

type FileTokens struct {
	Path      string
	TokenList []string
}

type BaseTokenizer struct {
	Path    string
	Scanner *bufio.Scanner
}

func NewBaseTokenizer(path string) BaseTokenizer {
	file, err := os.Open(path)
	common.Check(err)
	scanner := bufio.NewScanner(file)
	return BaseTokenizer{path, scanner}
}

func (base BaseTokenizer) Index(tokens chan FileTokens) {
	scanner := base.Scanner
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

	index := make(map[string]bool)
	reg, err := regexp.Compile("(^[^A-Za-z0-9]+)|([^A-Za-z0-9]+$)")
	common.Check(err)
	for scanner.Scan() {
		word := reg.ReplaceAllString(strings.ToLower(scanner.Text()), "")
		if len(word) > 0 && index[word] == false {
			index[word] = true
		}
	}

	words := make([]string, len(index))
	for key, _ := range index {
		words = append(words, key)
	}

	tokens <- FileTokens{base.Path, words}
}
