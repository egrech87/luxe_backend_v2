package glassEntity

import (
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/enums/glassType"
	"Quote/utils/mapUtils"
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/pkg/errors"
)

type GlassEntity struct {
	baseEntity.BaseEntity `json:"-" structs:"-"`
	FriendlyName          string  `json:"friendly_name" structs:"friendly_name"`
	Type                  *int    `json:"type" structs:"type"`
	UValue                float64 `json:"u_value" structs:"u_value"`
	PriceFactor           float64 `json:"price_factor" structs:"price_factor"`
	GlazingStructure      string  `json:"glazing_structure" structs:"glazing_structure"`
	Thickness             float64 `json:"thickness" structs:"thickness" gorm:"default:0"`
	DatasheetId           *uint64 `json:"datasheet_id" structs:"datasheet_id"` // The id of the datasheet file; null when there is none
}

func New() *GlassEntity {
	return new(GlassEntity)
}

func (e *GlassEntity) GetFriendlyName() string {
	return e.FriendlyName
}
func (e *GlassEntity) GetType() *int {
	return e.Type
}
func (e *GlassEntity) GetTypeName() *glassType.GlassType {
	key := e.GetType()
	if key == nil {
		return nil
	}
	return glassType.FromInt(*key)
}
func (e *GlassEntity) GetUValue() float64 {
	return e.UValue
}
func (e *GlassEntity) GetPriceFactor() float64 {
	return e.PriceFactor
}
func (e *GlassEntity) GetGlazingStructure() string {
	return e.GlazingStructure
}
func (e *GlassEntity) GetThickness() float64 {
	return e.Thickness
}
func (e *GlassEntity) GetDatasheetId() *uint64 {
	return e.DatasheetId
}

func (e *GlassEntity) SetFriendlyName(value string) {
	e.FriendlyName = value
}
func (e *GlassEntity) SetType(value *int) {
	e.Type = value
}
func (e *GlassEntity) SetTypeName(value *glassType.GlassType) {
	if value == nil {
		e.Type = nil
	} else {
		gt := int(*value)
		e.Type = &gt
	}
}
func (e *GlassEntity) SetUValue(value float64) {
	e.UValue = value
}
func (e *GlassEntity) SetPriceFactor(value float64) {
	e.PriceFactor = value
}
func (e *GlassEntity) SetGlazingStructure(value string) {
	e.GlazingStructure = value
}
func (e *GlassEntity) SetThickness(value float64) {
	e.Thickness = value
}
func (e *GlassEntity) SetDatasheetId(value *uint64) {
	e.DatasheetId = value
}

func (*GlassEntity) TableName() string {
	return TableName()
}
func TableName() string {
	return "glasses"
}
func (e *GlassEntity) FieldNames() []string {
	result := e.BaseEntity.FieldNames()
	return append(result, []string{"friendly_name", "type", "u_value", "price_factor", "glazing_structure", "thickness", "datasheet_id"}...)
}

func (e *GlassEntity) ToJson() map[string]any {
	json := structs.Map(e)

	for key, element := range e.BaseEntity.ToJson() {
		json[key] = element
	}

	return json
}

func (e *GlassEntity) FillFromJson(json map[string]any) error {
	if err := e.BaseEntity.FillFromJson(json); err != nil {
		return err
	}

	_, err := mapUtils.FieldAnyToType(json, "type", fmt.Sprintf("%T", e.Type))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "u_value", fmt.Sprintf("%T", e.UValue))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "price_factor", fmt.Sprintf("%T", e.PriceFactor))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "thickness", fmt.Sprintf("%T", e.Thickness))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "datasheet_id", fmt.Sprintf("%T", e.DatasheetId))
	if err != nil {
		return err
	}

	return mapUtils.ToStruct_usingJsonTags(json, e)
}

func (e *GlassEntity) FillFromDb(dbMap map[string]any) error {
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

func (e *GlassEntity) Clone() (baseEntity.BaseEntityI, error) {
	clone := New()
	err := clone.FillFromJson(e.ToJson())
	return clone, err
}

func (e *GlassEntity) Validate() error {
	// Validate the Type
	{
		key := e.GetType()
		if key != nil && !glassType.Exists(*key) {
			return errors.Errorf("Type '%d' is invalid", *key)
		}
	}

	return nil
}

func (e *GlassEntity) GetUpdates(oldItem baseEntity.BaseEntityI) map[string]any {
	oldGlass := oldItem.(*GlassEntity)
	updates := map[string]any{}

	if e.FriendlyName != oldGlass.FriendlyName {
		updates["friendly_name"] = e.FriendlyName
	}
	if e.Type != oldGlass.Type {
		updates["type"] = e.Type
	}
	if e.UValue != oldGlass.UValue {
		updates["u_value"] = e.UValue
	}
	if e.PriceFactor != oldGlass.PriceFactor {
		updates["price_factor"] = e.PriceFactor
	}
	if e.GlazingStructure != oldGlass.GlazingStructure {
		updates["glazing_structure"] = e.GlazingStructure
	}
	if e.Thickness != oldGlass.Thickness {
		updates["thickness"] = e.Thickness
	}
	if !equalIds(e.DatasheetId, oldGlass.DatasheetId) {
		updates["datasheet_id"] = e.DatasheetId
	}

	return updates
}

// equalIds compares two optional ids by value (nil only equals nil).
func equalIds(a, b *uint64) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
