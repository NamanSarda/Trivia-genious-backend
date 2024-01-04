package util

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

func CopyFileFromRequest(r *http.Request) ([]byte, error) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		return nil, err
	}

	file, _, err := r.FormFile("pdfFile")
	if err != nil {
		return nil, err
	}
	//	defer file.Close()

	// Read the file contents into a byte slice.
	fileContents, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return fileContents, nil
}

func ExtractTextFromPdf(pdfContents []byte) (string, error) {
	license.SetMeteredKey("1b46218f86a0fc3cdc8bb38c613b5f7567848e2b72787ac23e9dce15d77a61d3")

	pdfReader, err := model.NewPdfReader(bytes.NewReader(pdfContents))
	if err != nil {
		return "", err
	}
	//	defer pdfReader.Close()

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return "", err
	}

	var pdfText strings.Builder

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return "", err
		}

		extract, err := extractor.New(page)
		if err != nil {
			return "", err
		}

		text, err := extract.ExtractText()
		if err != nil {
			return "", err
		}

		pdfText.WriteString(text)
		pdfText.WriteString("\n")
	}

	return pdfText.String(), nil
}
