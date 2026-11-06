package components

import (
	"Quote/utils/mapUtils"

	"github.com/fatih/structs"
)

type Address struct {
	Block     string `json:"block" structs:"block"`
	Apartment string `json:"apartment" structs:"apartment"`
	Street    string `json:"street" structs:"street"`
	Locality  string `json:"locality" structs:"locality"`
	Country   string `json:"country" structs:"country"`
}

func (a *Address) GetBlock() string {
	return a.Block
}
func (a *Address) GetApartment() string {
	return a.Apartment
}
func (a *Address) GetStreet() string {
	return a.Street
}
func (a *Address) GetLocality() string {
	return a.Locality
}
func (a *Address) GetCountry() string {
	return a.Country
}
func (a *Address) SetBlock(value string) {
	a.Block = value
}
func (a *Address) SetApartment(value string) {
	a.Apartment = value
}
func (a *Address) SetStreet(value string) {
	a.Street = value
}
func (a *Address) SetLocality(value string) {
	a.Locality = value
}
func (a *Address) SetCountry(value string) {
	a.Country = value
}

func (a *Address) FieldNames() []string {
	return []string{"block", "apartment", "street", "locality", "country"}
}

func (a *Address) ToJson() map[string]any {
	return structs.Map(a)
}

func (a *Address) FillFromJson(json map[string]any) error {
	return mapUtils.ToStruct_usingJsonTags(json, a)
}

func (a *Address) Validate() error {
	return nil
}

func (a *Address) GetUpdates(oldAddress Address) map[string]any {
	updates := map[string]any{}

	if a.Block != oldAddress.Block {
		updates["block"] = a.Block
	}
	if a.Apartment != oldAddress.Apartment {
		updates["apartment"] = a.Apartment
	}
	if a.Street != oldAddress.Street {
		updates["street"] = a.Street
	}
	if a.Locality != oldAddress.Locality {
		updates["locality"] = a.Locality
	}
	if a.Country != oldAddress.Country {
		updates["country"] = a.Country
	}

	return updates
}
