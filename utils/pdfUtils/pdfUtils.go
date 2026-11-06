package pdfUtils

import (
	"bytes"
	"fmt"
	"io"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func MergePdfBytes(fileBytesList ...[]byte) ([]byte, error) {
	// Define the PDF processing configuration
	config := model.NewDefaultConfiguration()

	rsc := []io.ReadSeeker{}
	for _, fileBytes := range fileBytesList {
		rsc = append(rsc, bytes.NewReader(fileBytes))
	}

	var buffer bytes.Buffer // This is the writer

	// Merge the PDF files from memory (from bytes)
	err := api.MergeRaw(rsc, &buffer, false, config)
	if err != nil {
		fmt.Println("Error merging PDF files:", err)
		return []byte{}, err
	}

	return buffer.Bytes(), err
}
