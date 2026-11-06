package glassType

type GlassType int

const (
	Opaque GlassType = iota
	Clear
	Frosted
)

func List() []string {
	return []string{
		"Opaque",
		"Clear",
		"Frosted",
	}
}

func Exists(key int) bool {
	if key >= 0 && key < len(List()) {
		return true
	}

	return false
}

func FromString(name string) *GlassType {
	for i, v := range List() {
		if v == name {
			t := GlassType(i)
			return &t
		}
	}
	return nil
}

func FromInt(key int) *GlassType {
	if !Exists(key) {
		return nil
	}
	t := GlassType(key)
	return &t
}

func String(key int) string {
	t := FromInt(key)
	if t == nil {
		return ""
	}

	return t.String()
}

func (t GlassType) String() string {
	return List()[t]
}
