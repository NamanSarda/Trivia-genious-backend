package util

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

func CopyFileFromRequest(w http.ResponseWriter, r *http.Request) *os.File {
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	file, _, err := r.FormFile("pdfFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}
	defer file.Close()

	// Save the file to a temporary location.
	tempDir := os.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "uploaded*.pdf")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	defer os.Remove(tempFile.Name())
	return tempFile

}

func ExtractTextFromPdf(pdfPath string) string {
	// Set your license key if you have one.
	// license.SetLicenseFile(os.Getenv("LICENSE_KEY"))
	license.SetMeteredKey(os.Getenv("LICENSE_KEY"))
	// Specify the path to your PDF file.

	// Open the PDF file.
	file, err := os.Open(pdfPath)
	if err != nil {
		log.Fatalf("Error opening PDF file: %v", err)
	}
	defer file.Close()

	// Create a new PdfReader using the file handle.
	pdfReader, err := model.NewPdfReader(file)
	if err != nil {
		log.Fatalf("Error creating PDF reader: %v", err)
	}

	// Extract text from all pages.
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		log.Fatal("No Pages found by Reader")
	}

	pdfText := ""

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		// Get the PdfPage for the current page.
		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			log.Fatalf("Error getting page %d: %v", pageNum, err)
		}

		// Create a text extractor for the page.
		extract, err := extractor.New(page)
		if err != nil {
			log.Fatalf("Error creating text extractor for page %d: %v", pageNum, err)
		}

		// Extract text from the page.
		text, err := extract.ExtractText()
		if err != nil {
			log.Fatalf("Error extracting text from page %d: %v", pageNum, err)
		}

		// fmt.Printf("Page %d:\n%s\n", pageNum, text)
		pdfText += (text + "\n")

	}

	return pdfText
}
