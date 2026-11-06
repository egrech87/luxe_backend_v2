package pathUtils

import (
	"Quote/utils/stringUtils"
	"fmt"
	"os"
	"strings"
)

var rootPath string

// Set path to nil to have the path auto-generated. Note that this fails when the program is a windows service.
func Init(path *string) error {
	if path != nil {
		rootPath = *path
	} else {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}

		if strings.HasSuffix(dir, "cmd") {
			rootPath = stringUtils.BeforeFirstSubstring(dir, "\\cmd")
		} else {
			rootPath = dir
		}
	}

	return nil
}

func GetRootPath() string {
	return rootPath
}

func ToImagePath(imageName string) string {
	return fmt.Sprintf("%s\\images\\%s", rootPath, imageName)
}
