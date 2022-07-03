package readers

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ledongthuc/pdf"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/ryanmahan/gofind/common"
)

type PDFTokenizer struct {
	BaseTokenizer
}

func NewPDFTokenizer(path string) Tokenizer {
	file, _, err := pdf.Open(path)
	readseeker := io.ReadSeeker(file);

	defer file.Close()
	common.Check(err)

	pdfctx, err := pdfcpu.Read(readseeker, nil)
	num, str := pdfctx.Read.XRefStreamsString()
	fmt.Println(num, str)
	rdr, err := pdfctx.ExtractPageContent(0)
	
	common.Check(err)
	scanner := bufio.NewScanner(rdr)
	fmt.Println(scanner.Text())
	return PDFTokenizer{BaseTokenizer{path, scanner}}
}

func (pdfTokenizer PDFTokenizer) Index() []string {
	return pdfTokenizer.BaseTokenizer.Index()
}