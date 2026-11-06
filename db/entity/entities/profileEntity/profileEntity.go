package profileEntity

import (
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/supplierEntity"
	"Quote/utils/mapUtils"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/pkg/errors"
)

type ProfileEntity struct {
	baseEntity.BaseEntity `json:"-" structs:"-"`
	FriendlyName          string                         `json:"friendly_name" structs:"friendly_name"`
	Brand                 string                         `json:"brand" structs:"brand"`
	Serie                 string                         `json:"serie" structs:"serie"`
	Material              string                         `json:"material" structs:"material"`
	SystemThickness       float64                        `json:"system_thickness" structs:"system_thickness"`
	PriceMatrix           string                         `json:"price_matrix" structs:"price_matrix"` // json data
	LeftImageId           *uint64                        `json:"left_image_id" structs:"left_image_id"`
	RightImageId          *uint64                        `json:"right_image_id" structs:"right_image_id"`
	UValue                float64                        `json:"u_value" structs:"u_value"`
	DatasheetId           *uint64                        `json:"datasheet_id" structs:"datasheet_id"`
	SupplierID            *uint64                        `json:"supplier_id" structs:"supplier_id"`
	Supplier              *supplierEntity.SupplierEntity `json:"-" structs:"-"`
}

func New() *ProfileEntity {
	entity := new(ProfileEntity)
	entity.init()
	return entity
}

func (e *ProfileEntity) init() {
	e.PriceMatrix = "{}"
}

func (e *ProfileEntity) GetFriendlyName() string {
	return e.FriendlyName
}
func (e *ProfileEntity) GetBrand() string {
	return e.Brand
}
func (e *ProfileEntity) GetSerie() string {
	return e.Serie
}
func (e *ProfileEntity) GetMaterial() string {
	return e.Material
}
func (e *ProfileEntity) GetSystemThickness() float64 {
	return e.SystemThickness
}
func (e *ProfileEntity) GetPriceMatrix() string {
	return e.PriceMatrix
}
func (e *ProfileEntity) GetImageLeftId() *uint64 {
	return e.LeftImageId
}
func (e *ProfileEntity) GetImageRightId() *uint64 {
	return e.RightImageId
}
func (e *ProfileEntity) GetUValue() float64 {
	return e.UValue
}
func (e *ProfileEntity) GetDatasheetId() *uint64 {
	return e.DatasheetId
}
func (e *ProfileEntity) GetSupplierID() *uint64 {
	return e.SupplierID
}
func (e *ProfileEntity) SetFriendlyName(value string) {
	e.FriendlyName = value
}
func (e *ProfileEntity) SetBrand(value string) {
	e.Brand = value
}
func (e *ProfileEntity) SetSerie(value string) {
	e.Serie = value
}
func (e *ProfileEntity) SetMaterial(value string) {
	e.Material = value
}
func (e *ProfileEntity) SetSystemThickness(value float64) {
	e.SystemThickness = value
}
func (e *ProfileEntity) SetPriceMatrix(value string) {
	e.PriceMatrix = value
}
func (e *ProfileEntity) SetLeftImageId(value *uint64) {
	e.LeftImageId = value
}
func (e *ProfileEntity) SetRightImageId(value *uint64) {
	e.RightImageId = value
}
func (e *ProfileEntity) SetUValue(value float64) {
	e.UValue = value
}
func (e *ProfileEntity) SetDatasheetId(value *uint64) {
	e.DatasheetId = value
}
func (e *ProfileEntity) SetSupplierID(value *uint64) {
	e.SupplierID = value
}

func (*ProfileEntity) TableName() string {
	return TableName()
}
func TableName() string {
	return "profiles"
}

func (e *ProfileEntity) FieldNames() []string {
	result := e.BaseEntity.FieldNames()
	return append(result, []string{"friendly_name", "brand", "serie", "material", "system_thickness", "left_image_id", "right_image_id", "price_matrix", "u_value", "datasheet_id", "supplier_id"}...)
}

func (e *ProfileEntity) ToJson() map[string]any {
	json := structs.Map(e)

	// Alter price matrix to be a string
	json["price_matrix"] = string(e.PriceMatrix)

	for key, element := range e.BaseEntity.ToJson() {
		json[key] = element
	}

	if e.Supplier != nil {
		json["supplier"] = e.Supplier.ToJson()
	}

	return json
}

func (e *ProfileEntity) FillFromJson(json map[string]any) error {
	// Fix the price_matrix in case it is an empty string
	if json["price_matrix"] != nil && json["price_matrix"] == "" {
		delete(json, "price_matrix")
	}

	if err := e.BaseEntity.FillFromJson(json); err != nil {
		return err
	}

	_, err := mapUtils.FieldAnyToType(json, "price_matrix", fmt.Sprintf("%T", e.PriceMatrix))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "system_thickness", fmt.Sprintf("%T", e.SystemThickness))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "u_value", fmt.Sprintf("%T", e.UValue))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "left_image_id", fmt.Sprintf("%T", e.LeftImageId))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "right_image_id", fmt.Sprintf("%T", e.RightImageId))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "datasheet_id", fmt.Sprintf("%T", e.DatasheetId))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "supplier_id", fmt.Sprintf("%T", e.SupplierID))
	if err != nil {
		return err
	}

	return mapUtils.ToStruct_usingJsonTags(json, e)
}

func (e *ProfileEntity) FillFromDb(dbMap map[string]any) error {
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

	if e.SupplierID != nil {
		// This Profile has a supplier. Try to fill in the supplier's details
		itemMap := map[string]any{}
		supplier := supplierEntity.New()
		keyPrefix := supplier.TableName() + "_"
		for key, value := range dbMap {
			if strings.HasPrefix(key, keyPrefix) {
				newKey := strings.TrimPrefix(key, keyPrefix)
				itemMap[newKey] = value
			}
		}

		if len(itemMap) > 0 {
			if err := supplier.FillFromJson(itemMap); err != nil {
				return err
			}

			e.Supplier = supplier
		}
	}

	return nil
}

func (e *ProfileEntity) Clone() (baseEntity.BaseEntityI, error) {
	clone := New()
	err := clone.FillFromJson(e.ToJson())
	return clone, err
}

func (e *ProfileEntity) Validate() error {
	// Validate the Price Matrix
	{
		if !json.Valid([]byte(e.PriceMatrix)) {
			return errors.Errorf("Price Matrix is invalid")
		}
	}

	return nil
}

func (e *ProfileEntity) GetUpdates(oldItem baseEntity.BaseEntityI) map[string]any {
	oldProfile := oldItem.(*ProfileEntity)
	updates := map[string]any{}

	if e.FriendlyName != oldProfile.FriendlyName {
		updates["friendly_name"] = e.FriendlyName
	}
	if e.Brand != oldProfile.Brand {
		updates["brand"] = e.Brand
	}
	if e.Serie != oldProfile.Serie {
		updates["serie"] = e.Serie
	}
	if e.Material != oldProfile.Material {
		updates["material"] = e.Material
	}
	if e.SystemThickness != oldProfile.SystemThickness {
		updates["system_thickness"] = e.SystemThickness
	}
	if e.LeftImageId != oldProfile.LeftImageId {
		updates["left_image_id"] = e.LeftImageId
	}
	if e.RightImageId != oldProfile.RightImageId {
		updates["right_image_id"] = e.RightImageId
	}
	if e.UValue != oldProfile.UValue {
		updates["u_value"] = e.UValue
	}
	if e.PriceMatrix != oldProfile.PriceMatrix {
		updates["price_matrix"] = e.PriceMatrix
	}
	if e.DatasheetId != oldProfile.DatasheetId {
		updates["datasheet_id"] = e.DatasheetId
	}
	if e.SupplierID != oldProfile.SupplierID {
		updates["supplier_id"] = e.SupplierID
	}

	return updates
}
