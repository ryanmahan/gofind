package readers

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/ryanmahan/gofind/common"
)

type Tokenizer interface {
	Index() []string
}

type BaseTokenizer struct {
	Path string
	Scanner *bufio.Scanner
}

func NewBaseTokenizer(path string) Tokenizer {
	file, err := os.Open(path)
	common.Check(err)
	scanner := bufio.NewScanner(file)
	return BaseTokenizer{path, scanner}
}

func (base BaseTokenizer) Index() []string {
	scanner := base.Scanner
	scanner.Split(func(data []byte, eof bool) (advance int, token []byte, err error) {
		if eof {
			return len(data), data, bufio.ErrFinalToken
		}
		for index, token := range data {
			fmt.Println(string(token))
			if index > 0 && (token == '\n' || token == ' ') {
				return index + 1, data[0:index], nil
			}
		}
		return 0, nil, nil
	})

	var index []string;
	reg, err := regexp.Compile("(^[^A-Za-z0-9]+)|([^A-Za-z0-9]+$)")
	common.Check(err)
	fmt.Println("here!")
	for scanner.Scan() {
		word := reg.ReplaceAllString(strings.ToLower(scanner.Text()), "")
		fmt.Println(word)
		if len(word) > 0 {
			index = append(index, word)
		}
	}
	return index
}