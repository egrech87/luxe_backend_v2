package fileTypes

type FileType int

const (
	Png FileType = iota
	Jpg
	Gif
	Pdf
	WordDoc
	Excel
)

func List() []string {
	return []string{
		"Png",
		"Jpg",
		"Gif",
		"Pdf",
		"WordDoc",
		"Excel",
	}
}

func Exists(key int) bool {
	if key >= 0 && key < len(List()) {
		return true
	}

	return false
}

func FromString(name string) *FileType {
	for i, v := range List() {
		if v == name {
			t := FileType(i)
			return &t
		}
	}
	return nil
}

func FromInt(key int) *FileType {
	if !Exists(key) {
		return nil
	}
	t := FileType(key)
	return &t
}

func String(key int) string {
	t := FromInt(key)
	if t == nil {
		return ""
	}

	return t.String()
}

func (t FileType) String() string {
	return List()[t]
}

func (t FileType) IsImage() bool {
	return t == Png || t == Jpg || t == Gif
}
