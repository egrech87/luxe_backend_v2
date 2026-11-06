package clientEntity

import (
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/components"
	"Quote/utils/mapUtils"
	"strings"

	"github.com/fatih/structs"
)

type ClientEntity struct {
	baseEntity.BaseEntity `json:"-" structs:"-"`
	components.Fullname   `json:"-" structs:"-"`
	Address               string `json:"address" structs:"address"`
	Email                 string `json:"email" structs:"email"`
	ContactNumber         string `json:"contact_number" structs:"contact_number"`
}

func New() *ClientEntity {
	return new(ClientEntity)
}

func (e *ClientEntity) GetAddress() string {
	return e.Address
}
func (e *ClientEntity) GetEmail() string {
	return e.Email
}
func (e *ClientEntity) GetContactNumber() string {
	return e.ContactNumber
}
func (e *ClientEntity) SetAddress(value string) {
	e.Address = value
}
func (e *ClientEntity) SetEmail(value string) {
	e.Email = value
}
func (e *ClientEntity) SetContactNumber(value string) {
	e.ContactNumber = value
}

func (*ClientEntity) TableName() string {
	return TableName()
}
func TableName() string {
	return "clients"
}
func (e *ClientEntity) FieldNames() []string {
	result := e.BaseEntity.FieldNames()
	return append(result, []string{"name", "surname", "address", "email", "contact_number"}...)
}

func (e *ClientEntity) ToJson() map[string]any {
	json := structs.Map(e)

	for key, element := range e.BaseEntity.ToJson() {
		json[key] = element
	}
	for key, element := range e.Fullname.ToJson() {
		json[key] = element
	}

	return json
}

func (e *ClientEntity) FillFromJson(json map[string]any) error {
	if err := e.BaseEntity.FillFromJson(json); err != nil {
		return err
	}

	if err := e.Fullname.FillFromJson(json); err != nil {
		return err
	}

	return mapUtils.ToStruct_usingJsonTags(json, e)
}

func (e *ClientEntity) FillFromDb(dbMap map[string]any) error {
	itemMap := map[string]any{}
	keyPrefix := e.TableName() + "_"
	for key, value := range dbMap {
		if strings.HasPrefix(key, keyPrefix) {
			newKey := strings.TrimPrefix(key, keyPrefix)
			itemMap[newKey] = value
		}
	}

	if err := e.FillFromJson(itemMap); err != nil {
		return err
	}

	return nil
}

func (e *ClientEntity) Clone() (baseEntity.BaseEntityI, error) {
	clone := New()
	err := clone.FillFromJson(e.ToJson())
	return clone, err
}

func (e *ClientEntity) Validate() error {
	if err := e.Fullname.Validate(); err != nil {
		return err
	}

	return nil
}

func (e *ClientEntity) GetUpdates(oldItem baseEntity.BaseEntityI) map[string]any {
	oldClient := oldItem.(*ClientEntity)
	updates := map[string]any{}

	for key, element := range e.Fullname.GetUpdates(oldClient.Fullname) {
		updates[key] = element
	}

	if e.Address != oldClient.Address {
		updates["address"] = e.Address
	}
	if e.Email != oldClient.Email {
		updates["email"] = e.Email
	}
	if e.ContactNumber != oldClient.ContactNumber {
		updates["contact_number"] = e.ContactNumber
	}

	return updates
}
