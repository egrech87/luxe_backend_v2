package pdfUtils

import (
	"bytes"
	"fmt"
	"io"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
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

// AddPageNumbers stamps "Page n of N" at the bottom centre of every page of the given PDF.
func AddPageNumbers(pdfBytes []byte) ([]byte, error) {
	config := model.NewDefaultConfiguration()

	// %p is the current page number and %P the total number of pages
	wm, err := api.TextWatermark("Page %p of %P", "pos:bc, off:0 8, scale:1 abs, points:9, rot:0, fillc:#444444", true, false, types.POINTS)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	if err := api.AddWatermarks(bytes.NewReader(pdfBytes), &buffer, nil, wm, config); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
