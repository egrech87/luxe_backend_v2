package imageUtils

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

func FromMultipartFileHeader(fileHeaders []*multipart.FileHeader) ([][]byte, error) {
	imagesBytes := [][]byte{}

	for _, file := range fileHeaders {
		openedFile, err := file.Open()
		{
			if err != nil {
				return imagesBytes, errors.Errorf("failed to open the file: %s", err.Error())
			}
			defer openedFile.Close()
		}

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, openedFile); err != nil {
			return imagesBytes, errors.Errorf("failed to read bytes from the file: %s", err.Error())
		}

		imagesBytes = append(imagesBytes, buf.Bytes())
	}

	return imagesBytes, nil
}

func ResizeTo(img []byte) ([]byte, error) {
	imge, _, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		return []byte{}, err
	}

	resizedImage := resize.Resize(800, 0, imge, resize.Lanczos2)

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, resizedImage, nil)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func FromFilePath(imagePath string) ([]byte, error) {
	return os.ReadFile(imagePath)
}
