package apertureType

type ApertureType int

const (
	Door ApertureType = iota
	Window
)

func List() []string {
	return []string{
		"Door",
		"Window",
	}
}

func Exists(key int) bool {
	if key >= 0 && key < len(List()) {
		return true
	}

	return false
}

func FromString(name string) *ApertureType {
	for i, v := range List() {
		if v == name {
			t := ApertureType(i)
			return &t
		}
	}
	return nil
}

func FromInt(key int) *ApertureType {
	if !Exists(key) {
		return nil
	}
	t := ApertureType(key)
	return &t
}

func String(key int) string {
	t := FromInt(key)
	if t == nil {
		return ""
	}

	return t.String()
}

func (t ApertureType) String() string {
	return List()[t]
}
