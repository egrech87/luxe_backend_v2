package fileUtils

import (
	"bytes"
	"io"
)

func File_GetContent(reader io.Reader) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, reader); err != nil {
		return []byte{}, err
	}

	return buf.Bytes(), nil
}
