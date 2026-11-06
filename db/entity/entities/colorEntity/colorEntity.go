package colorEntity

import (
	"Quote/db/entity/baseEntity"
	"Quote/utils/mapUtils"
	"fmt"
	"strings"

	"github.com/fatih/structs"
)

type ColorEntity struct {
	baseEntity.BaseEntity `json:"-" structs:"-"`
	FriendlyName          string  `json:"friendly_name" structs:"friendly_name"`
	PriceFactor           float64 `json:"price_factor" structs:"price_factor"`
	Code                  string  `json:"code" structs:"code"`
}

func New() *ColorEntity {
	return new(ColorEntity)
}

func (e *ColorEntity) GetFriendlyName() string {
	return e.FriendlyName
}
func (e *ColorEntity) GetPriceFactor() float64 {
	return e.PriceFactor
}
func (e *ColorEntity) GetCode() string {
	return e.Code
}
func (e *ColorEntity) SetFriendlyName(value string) {
	e.FriendlyName = value
}
func (e *ColorEntity) SetPriceFactor(value float64) {
	e.PriceFactor = value
}
func (e *ColorEntity) SetCode(value string) {
	e.Code = value
}

func (*ColorEntity) TableName() string {
	return TableName()
}
func TableName() string {
	return "colors"
}
func (e *ColorEntity) FieldNames() []string {
	result := e.BaseEntity.FieldNames()
	return append(result, []string{"friendly_name", "price_factor", "code"}...)
}

func (e *ColorEntity) ToJson() map[string]any {
	json := structs.Map(e)

	for key, element := range e.BaseEntity.ToJson() {
		json[key] = element
	}

	return json
}

func (e *ColorEntity) FillFromJson(json map[string]any) error {
	if err := e.BaseEntity.FillFromJson(json); err != nil {
		return err
	}

	_, err := mapUtils.FieldAnyToType(json, "price_factor", fmt.Sprintf("%T", e.PriceFactor))
	if err != nil {
		return err
	}

	return mapUtils.ToStruct_usingJsonTags(json, e)
}

func (e *ColorEntity) FillFromDb(dbMap map[string]any) error {
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

func (e *ColorEntity) Clone() (baseEntity.BaseEntityI, error) {
	clone := New()
	err := clone.FillFromJson(e.ToJson())
	return clone, err
}

func (e *ColorEntity) Validate() error {
	return nil
}

func (e *ColorEntity) GetUpdates(oldItem baseEntity.BaseEntityI) map[string]any {
	oldColor := oldItem.(*ColorEntity)
	updates := map[string]any{}

	if e.FriendlyName != oldColor.FriendlyName {
		updates["friendly_name"] = e.FriendlyName
	}
	if e.PriceFactor != oldColor.PriceFactor {
		updates["price_factor"] = e.PriceFactor
	}
	if e.Code != oldColor.Code {
		updates["code"] = e.Code
	}

	return updates
}
