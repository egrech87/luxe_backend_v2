package quoteEntity

import (
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/components"
	"Quote/db/entity/entities/clientEntity"
	"Quote/db/entity/entities/userEntity"
	"Quote/utils/mapUtils"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/structs"
)

type QuoteEntity struct {
	baseEntity.BaseEntity `json:"-" structs:"-"`
	components.Address    `json:"-" structs:"-"`
	Number                int                        `json:"number" structs:"number"` // The quote number
	FriendlyName          string                     `json:"friendly_name" structs:"friendly_name"`
	Date                  *time.Time                 `json:"date" structs:"date"`
	ContingencyPerc       float64                    `json:"contingency_perc" structs:"contingency_perc"`
	ShippingPerUnit       float64                    `json:"shipping_per_unit" structs:"shipping_per_unit"`         // per square metre
	InstallationPerUnit   float64                    `json:"installation_per_unit" structs:"installation_per_unit"` // per square metre
	DiscountPerc          float64                    `json:"discount_perc" structs:"discount_perc"`
	GrossProfit           float64                    `json:"gross_profit" structs:"gross_profit"`
	IsGrossProfitCustom   bool                       `json:"is_gross_profit_custom" structs:"is_gross_profit_custom" gorm:"default:false"`
	ProfitMargin          float64                    `json:"profit_margin" structs:"profit_margin"`
	Vat                   float64                    `json:"vat" structs:"vat"`
	IsVatCustom           bool                       `json:"is_vat_custom" structs:"is_vat_custom" gorm:"default:false"`
	SubTotal              float64                    `json:"sub_total" structs:"sub_total"`
	IsSubTotalCustom      bool                       `json:"is_sub_total_custom" structs:"is_sub_total_custom" gorm:"default:false"`
	Total                 float64                    `json:"total" structs:"total"`
	IsTotalCustom         bool                       `json:"is_total_custom" structs:"is_total_custom" gorm:"default:false"`
	Discount              float64                    `json:"discount" structs:"discount"`
	IsDiscountCustom      bool                       `json:"is_discount_custom" structs:"is_discount_custom" gorm:"default:false"`
	TotalArea             float64                    `json:"total_area" structs:"total_area" gorm:"-"`
	TotalCost             float64                    `json:"total_cost" structs:"total_cost" gorm:"-"`
	TotalApertures        int                        `json:"total_apertures" structs:"total_apertures" gorm:"-"`
	NetsQuantity          int                        `json:"nets_quantity" structs:"nets_quantity" gorm:"-"`
	ClientID              *uint64                    `json:"client_id" structs:"client_id"`
	Client                *clientEntity.ClientEntity `json:"-" structs:"-"`
	CreateUser            *userEntity.UserEntity     `json:"-" structs:"-" gorm:"foreignKey:created_by"` // The user who created the quote
	UpdateUser            *userEntity.UserEntity     `json:"-" structs:"-" gorm:"foreignKey:updated_by"` // The user who last updated the quote
}

func New() *QuoteEntity {
	entity := new(QuoteEntity)
	entity.init()
	return entity
}

func (e *QuoteEntity) init() {
	// Intentionally left empty
}

func (e *QuoteEntity) GetNumber() int {
	return e.Number
}
func (e *QuoteEntity) GetFriendlyName() string {
	return e.FriendlyName
}
func (e *QuoteEntity) GetDate() *time.Time {
	return e.Date
}
func (e *QuoteEntity) GetContingencyPerc() float64 {
	return e.ContingencyPerc
}
func (e *QuoteEntity) GetShippingPerUnit() float64 {
	return e.ShippingPerUnit
}
func (e *QuoteEntity) GetInstallationPerUnit() float64 {
	return e.InstallationPerUnit
}
func (e *QuoteEntity) GetDiscountPerc() float64 {
	return e.DiscountPerc
}
func (e *QuoteEntity) GetGrossProfit() float64 {
	return e.GrossProfit
}
func (e *QuoteEntity) GetIsGrossProfitCustom() bool {
	return e.IsGrossProfitCustom
}
func (e *QuoteEntity) GetProfitMargin() float64 {
	return e.ProfitMargin
}
func (e *QuoteEntity) GetVAT() float64 {
	return e.Vat
}
func (e *QuoteEntity) GetIsVatCustom() bool {
	return e.IsVatCustom
}
func (e *QuoteEntity) GetSubTotal() float64 {
	return e.SubTotal
}
func (e *QuoteEntity) GetIsSubTotalCustom() bool {
	return e.IsSubTotalCustom
}
func (e *QuoteEntity) GetTotal() float64 {
	return e.Total
}
func (e *QuoteEntity) GetIsTotalCustom() bool {
	return e.IsTotalCustom
}
func (e *QuoteEntity) GetDiscount() float64 {
	return e.Discount
}
func (e *QuoteEntity) GetIsDiscountCustom() bool {
	return e.IsDiscountCustom
}
func (e *QuoteEntity) GetTotalArea() float64 {
	return e.TotalArea
}
func (e *QuoteEntity) GetTotalCost() float64 {
	return e.TotalCost
}
func (e *QuoteEntity) GetTotalApertures() int {
	return e.TotalApertures
}
func (e *QuoteEntity) GetNetsQuantity() int {
	return e.NetsQuantity
}
func (e *QuoteEntity) GetClientID() *uint64 {
	return e.ClientID
}
func (e *QuoteEntity) SetNumber(value int) {
	e.Number = value
}
func (e *QuoteEntity) SetFriendlyName(value string) {
	e.FriendlyName = value
}
func (e *QuoteEntity) SetDate(value *time.Time) {
	e.Date = value
}
func (e *QuoteEntity) SetContingencyPerc(value float64) {
	e.ContingencyPerc = value
}
func (e *QuoteEntity) SetShippingPerUnit(value float64) {
	e.ShippingPerUnit = value
}
func (e *QuoteEntity) SetInstallationPerUnit(value float64) {
	e.InstallationPerUnit = value
}
func (e *QuoteEntity) SetDiscountPerc(value float64) {
	e.DiscountPerc = value
}
func (e *QuoteEntity) SetGrossProfit(value float64) {
	e.GrossProfit = value
}
func (e *QuoteEntity) SetIsGrossProfitCustom(value bool) {
	e.IsGrossProfitCustom = value
}
func (e *QuoteEntity) SetProfitMargin(value float64) {
	e.ProfitMargin = value
}
func (e *QuoteEntity) SetVAT(value float64) {
	e.Vat = value
}
func (e *QuoteEntity) SetIsVatCustom(value bool) {
	e.IsVatCustom = value
}
func (e *QuoteEntity) SetSubTotal(value float64) {
	e.SubTotal = value
}
func (e *QuoteEntity) SetIsSubTotalCustom(value bool) {
	e.IsSubTotalCustom = value
}
func (e *QuoteEntity) SetTotal(value float64) {
	e.Total = value
}
func (e *QuoteEntity) SetIsTotalCustom(value bool) {
	e.IsTotalCustom = value
}
func (e *QuoteEntity) SetDiscount(value float64) {
	e.Discount = value
}
func (e *QuoteEntity) SetIsDiscountCustom(value bool) {
	e.IsDiscountCustom = value
}
func (e *QuoteEntity) SetClientID(value *uint64) {
	e.ClientID = value
}

func (*QuoteEntity) TableName() string {
	return TableName()
}
func TableName() string {
	return "quotes"
}

func (e *QuoteEntity) FieldNames() []string {
	result := e.BaseEntity.FieldNames()
	result = append(result, e.Address.FieldNames()...)
	return append(result, []string{"number", "friendly_name", "date", "contingency_perc", "shipping_per_unit", "installation_per_unit", "discount_perc", "is_gross_profit_custom", "profit_margin", "is_vat_custom", "is_sub_total_custom", "is_total_custom", "is_discount_custom", "client_id"}...)
}

func (e *QuoteEntity) ToJson() map[string]any {
	json := structs.Map(e)

	for key, element := range e.BaseEntity.ToJson() {
		json[key] = element
	}

	for key, element := range e.Address.ToJson() {
		json[key] = element
	}

	if e.Client != nil {
		json["client"] = e.Client.ToJson()
	}

	if e.CreateUser != nil {
		json["create_user"] = e.CreateUser.ToJson()
	}
	if e.UpdateUser != nil {
		json["update_user"] = e.UpdateUser.ToJson()
	}

	return json
}

func (e *QuoteEntity) FillFromJson(json map[string]any) error {
	if err := e.BaseEntity.FillFromJson(json); err != nil {
		return err
	}

	if err := e.Address.FillFromJson(json); err != nil {
		return err
	}

	_, err := mapUtils.FieldAnyToType(json, "number", fmt.Sprintf("%T", e.Number))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "date", fmt.Sprintf("%T", e.Date))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "contingency_perc", fmt.Sprintf("%T", e.ContingencyPerc))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "shipping_per_unit", fmt.Sprintf("%T", e.ShippingPerUnit))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "installation_per_unit", fmt.Sprintf("%T", e.InstallationPerUnit))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "discount_perc", fmt.Sprintf("%T", e.DiscountPerc))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "gross_profit", fmt.Sprintf("%T", e.GrossProfit))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "profit_margin", fmt.Sprintf("%T", e.ProfitMargin))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "vat", fmt.Sprintf("%T", e.Vat))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "sub_total", fmt.Sprintf("%T", e.SubTotal))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "total", fmt.Sprintf("%T", e.Total))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "discount", fmt.Sprintf("%T", e.Discount))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "total_area", fmt.Sprintf("%T", e.TotalArea))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "total_cost", fmt.Sprintf("%T", e.TotalCost))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "total_apertures", fmt.Sprintf("%T", e.TotalApertures))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "nets_quantity", fmt.Sprintf("%T", e.NetsQuantity))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "client_id", fmt.Sprintf("%T", e.ClientID))
	if err != nil {
		return err
	}

	return mapUtils.ToStruct_usingJsonTags(json, e)
}

func (e *QuoteEntity) FillFromDb(dbMap map[string]any) error {
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

	if e.ClientID != nil {
		// This Quote has a client. Try to fill in the client's details
		itemMap := map[string]any{}
		client := clientEntity.New()
		keyPrefix := client.TableName() + "_"
		for key, value := range dbMap {
			if strings.HasPrefix(key, keyPrefix) {
				newKey := strings.TrimPrefix(key, keyPrefix)
				itemMap[newKey] = value
			}
		}

		if len(itemMap) > 0 {
			if err := client.FillFromJson(itemMap); err != nil {
				return err
			}

			e.Client = client
		}
	}

	{
		// This Quote has a create user. Try to fill in the user's details
		itemMap := map[string]any{}
		user := userEntity.New()
		keyPrefix := "createuser_"
		for key, value := range dbMap {
			if strings.HasPrefix(key, keyPrefix) {
				newKey := strings.TrimPrefix(key, keyPrefix)
				itemMap[newKey] = value
			}
		}

		if len(itemMap) > 0 {
			if err := user.FillFromJson(itemMap); err != nil {
				return err
			}

			e.CreateUser = user
		}
	}

	{
		// This Quote has an update user. Try to fill in the user's details
		itemMap := map[string]any{}
		user := userEntity.New()
		keyPrefix := "updateuser_"
		for key, value := range dbMap {
			if strings.HasPrefix(key, keyPrefix) {
				newKey := strings.TrimPrefix(key, keyPrefix)
				itemMap[newKey] = value
			}
		}

		if len(itemMap) > 0 {
			if err := user.FillFromJson(itemMap); err != nil {
				return err
			}

			e.UpdateUser = user
		}
	}

	return nil
}

func (e *QuoteEntity) Clone() (baseEntity.BaseEntityI, error) {
	clone := New()
	err := clone.FillFromJson(e.ToJson())
	return clone, err
}

func (e *QuoteEntity) Validate() error {
	return nil
}

func (e *QuoteEntity) GetUpdates(oldItem baseEntity.BaseEntityI) map[string]any {
	oldQuote := oldItem.(*QuoteEntity)
	updates := map[string]any{}

	for key, element := range e.Address.GetUpdates(oldQuote.Address) {
		updates[key] = element
	}

	if e.Number != oldQuote.Number {
		updates["number"] = e.Number
	}
	if e.FriendlyName != oldQuote.FriendlyName {
		updates["friendly_name"] = e.FriendlyName
	}
	if e.Date != oldQuote.Date {
		updates["date"] = e.Date
	}
	if e.ContingencyPerc != oldQuote.ContingencyPerc {
		updates["contingency_perc"] = e.ContingencyPerc
	}
	if e.ShippingPerUnit != oldQuote.ShippingPerUnit {
		updates["shipping_per_unit"] = e.ShippingPerUnit
	}
	if e.InstallationPerUnit != oldQuote.InstallationPerUnit {
		updates["installation_per_unit"] = e.InstallationPerUnit
	}
	if e.DiscountPerc != oldQuote.DiscountPerc {
		updates["discount_perc"] = e.DiscountPerc
	}
	if e.GrossProfit != oldQuote.GrossProfit {
		updates["gross_profit"] = e.GrossProfit
	}
	if e.IsGrossProfitCustom != oldQuote.IsGrossProfitCustom {
		updates["is_gross_profit_custom"] = e.IsGrossProfitCustom
	}
	if e.ProfitMargin != oldQuote.ProfitMargin {
		updates["profit_margin"] = e.ProfitMargin
	}
	if e.Vat != oldQuote.Vat {
		updates["vat"] = e.Vat
	}
	if e.IsVatCustom != oldQuote.IsVatCustom {
		updates["is_vat_custom"] = e.IsVatCustom
	}
	if e.SubTotal != oldQuote.SubTotal {
		updates["sub_total"] = e.SubTotal
	}
	if e.IsSubTotalCustom != oldQuote.IsSubTotalCustom {
		updates["is_sub_total_custom"] = e.IsSubTotalCustom
	}
	if e.Total != oldQuote.Total {
		updates["total"] = e.Total
	}
	if e.IsTotalCustom != oldQuote.IsTotalCustom {
		updates["is_total_custom"] = e.IsTotalCustom
	}
	if e.Discount != oldQuote.Discount {
		updates["discount"] = e.Discount
	}
	if e.IsDiscountCustom != oldQuote.IsDiscountCustom {
		updates["is_discount_custom"] = e.IsDiscountCustom
	}
	if e.ClientID != oldQuote.ClientID {
		updates["client_id"] = e.ClientID
	}

	return updates
}

func (e *QuoteEntity) HasNets() bool {
	if e.NetsQuantity > 0 {
		return true
	}

	return false
}
