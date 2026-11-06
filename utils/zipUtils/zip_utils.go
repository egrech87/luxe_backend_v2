package zipUtils

import (
	"archive/zip"
	"bytes"
	"fmt"
)

// Receives a list of files
func Zip_Compress(filenames []string, uncompressed [][]byte) ([]byte, error) {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	for i, file := range uncompressed {
		f, err := w.Create(filenames[i])
		if err != nil {
			return nil, err
		}

		_, err = f.Write(file)
		if err != nil {
			return nil, err
		}
	}

	// Make sure to check the error on Close.
	err := w.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Decompresses a zip file. Shall return a list of files from the zip file.
// Note: Function still in progress!!
func Zip_DecompressFile(filename string) error {
	zipReader, err := zip.OpenReader(filename)
	if err != nil {
		panic(err)
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fmt.Println(f.Name)
	}
	return nil
}

// Decompresses a zip file. Returns a list of files from the zip file.
func Zip_DecompressBytes(zippedBytes []byte) ([]*zip.File, error) {
	reader := bytes.NewReader(zippedBytes)
	zipReader, err := zip.NewReader(reader, int64(len(zippedBytes)))
	if err != nil {
		return nil, err
	}

	return zipReader.File, nil
}
