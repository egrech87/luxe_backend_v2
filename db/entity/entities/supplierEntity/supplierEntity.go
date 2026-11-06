package supplierEntity

import (
	"Quote/db/entity/baseEntity"
	"Quote/utils/mapUtils"
	"strings"

	"github.com/fatih/structs"
)

type SupplierEntity struct {
	baseEntity.BaseEntity `json:"-" structs:"-"`
	FriendlyName          string `json:"friendly_name" structs:"friendly_name"`
	Name                  string `json:"name" structs:"name"`
	Address               string `json:"address" structs:"address"`
	ContactNumber         string `json:"contact_number" structs:"contact_number"`
	Email                 string `json:"email" structs:"email"`
	HasPriceMatrix        bool   `json:"has_price_matrix" structs:"has_price_matrix"`
}

func New() *SupplierEntity {
	entity := new(SupplierEntity)
	entity.init()
	return entity
}

func (e *SupplierEntity) init() {
}

func (e *SupplierEntity) GetFriendlyName() string {
	return e.FriendlyName
}
func (e *SupplierEntity) GetName() string {
	return e.Name
}
func (e *SupplierEntity) GetAddress() string {
	return e.Address
}
func (e *SupplierEntity) GetContactNumber() string {
	return e.ContactNumber
}
func (e *SupplierEntity) GetEmail() string {
	return e.Email
}
func (e *SupplierEntity) GetHasPriceMatrix() bool {
	return e.HasPriceMatrix
}
func (e *SupplierEntity) SetFriendlyName(value string) {
	e.FriendlyName = value
}
func (e *SupplierEntity) SetName(value string) {
	e.Name = value
}
func (e *SupplierEntity) SetAddress(value string) {
	e.Address = value
}
func (e *SupplierEntity) SetContactNumber(value string) {
	e.ContactNumber = value
}
func (e *SupplierEntity) SetEmail(value string) {
	e.Email = value
}
func (e *SupplierEntity) SetHasPriceMatrix(value bool) {
	e.HasPriceMatrix = value
}

func (*SupplierEntity) TableName() string {
	return TableName()
}
func TableName() string {
	return "suppliers"
}
func (e *SupplierEntity) FieldNames() []string {
	result := e.BaseEntity.FieldNames()
	return append(result, []string{"friendly_name", "name", "address", "contact_number", "email", "has_price_matrix"}...)
}

func (e *SupplierEntity) ToJson() map[string]any {
	json := structs.Map(e)

	for key, element := range e.BaseEntity.ToJson() {
		json[key] = element
	}

	return json
}

func (e *SupplierEntity) FillFromJson(json map[string]any) error {
	if err := e.BaseEntity.FillFromJson(json); err != nil {
		return err
	}

	return mapUtils.ToStruct_usingJsonTags(json, e)
}

func (e *SupplierEntity) FillFromDb(dbMap map[string]any) error {
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

func (e *SupplierEntity) Clone() (baseEntity.BaseEntityI, error) {
	clone := New()
	err := clone.FillFromJson(e.ToJson())
	return clone, err
}

func (e *SupplierEntity) Validate() error {
	return nil
}

func (e *SupplierEntity) GetUpdates(oldItem baseEntity.BaseEntityI) map[string]any {
	oldSupplier := oldItem.(*SupplierEntity)
	updates := map[string]any{}

	if e.FriendlyName != oldSupplier.FriendlyName {
		updates["friendly_name"] = e.FriendlyName
	}
	if e.Name != oldSupplier.Name {
		updates["name"] = e.Name
	}
	if e.Address != oldSupplier.Address {
		updates["address"] = e.Address
	}
	if e.ContactNumber != oldSupplier.ContactNumber {
		updates["contact_number"] = e.ContactNumber
	}
	if e.Email != oldSupplier.Email {
		updates["email"] = e.Email
	}
	if e.HasPriceMatrix != oldSupplier.HasPriceMatrix {
		updates["has_price_matrix"] = e.HasPriceMatrix
	}

	return updates
}
