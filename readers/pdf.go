package readers

import (
	"os"
	"os/exec"

	"github.com/ryanmahan/gofind/common"
)

type PDFTokenizer struct {
	BaseTokenizer
}

func NewPDFTokenizer(path string) Tokenizer {

	convertCmd := exec.Command("pdftotext", "-q", "-nopgbrk", "-enc", "UTF-8", "-eol", "unix", path, "-")
	outfile, err := os.Create(path + ".txt")
	convertCmd.Stdout = outfile
	convertCmd.Start()
	convertCmd.Wait()
	common.Check(err)
	return NewBaseTokenizer(path + ".txt")
}

func (pdfTokenizer PDFTokenizer) Index(channel chan FileTokens) {
	pdfTokenizer.BaseTokenizer.Index(channel)
}
