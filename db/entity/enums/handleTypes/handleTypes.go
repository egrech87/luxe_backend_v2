package handleType

type HandleType int

const (
	Left HandleType = iota
	Right
)

func List() []string {
	return []string{
		"Left",
		"Right",
	}
}

func Exists(key int) bool {
	if key >= 0 && key < len(List()) {
		return true
	}

	return false
}

func FromString(name string) *HandleType {
	for i, v := range List() {
		if v == name {
			t := HandleType(i)
			return &t
		}
	}
	return nil
}

func FromInt(key int) *HandleType {
	if !Exists(key) {
		return nil
	}
	t := HandleType(key)
	return &t
}

func String(key int) string {
	t := FromInt(key)
	if t == nil {
		return ""
	}

	return t.String()
}

func (t HandleType) String() string {
	return List()[t]
}
