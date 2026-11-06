package components

import (
	"Quote/utils/mapUtils"
	"strings"

	"github.com/fatih/structs"
)

type Fullname struct {
	Name    string `json:"name" structs:"name"`
	Surname string `json:"surname" structs:"surname"`
}

func (f *Fullname) GetName() string {
	return f.Name
}
func (f *Fullname) GetSurname() string {
	return f.Surname
}
func (f *Fullname) SetName(value string) {
	f.Name = value
}
func (f *Fullname) SetSurname(value string) {
	f.Surname = value
}

func (f *Fullname) GetFullName() string {
	return strings.TrimSpace(f.Name + " " + f.Surname)
}

func (f *Fullname) ToJson() map[string]any {
	return structs.Map(f)
}

func (f *Fullname) FillFromJson(json map[string]any) error {
	return mapUtils.ToStruct_usingJsonTags(json, f)
}

func (f *Fullname) Validate() error {
	return nil
}

func (f *Fullname) GetUpdates(oldFullname Fullname) map[string]any {
	updates := map[string]any{}

	if f.Name != oldFullname.Name {
		updates["name"] = f.Name
	}
	if f.Surname != oldFullname.Surname {
		updates["surname"] = f.Surname
	}

	return updates
}
