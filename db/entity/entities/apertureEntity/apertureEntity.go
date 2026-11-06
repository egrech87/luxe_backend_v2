package apertureEntity

import (
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/colorEntity"
	"Quote/db/entity/entities/fileEntity"
	"Quote/db/entity/entities/glassEntity"
	"Quote/db/entity/entities/profileEntity"
	"Quote/db/entity/entities/quoteEntity"
	"Quote/db/entity/entities/userEntity"
	"Quote/db/entity/enums/apertureType"
	handleType "Quote/db/entity/enums/handleTypes"
	"Quote/utils/int64Array"
	"Quote/utils/mapUtils"
	"fmt"
	"slices"
	"strings"

	"github.com/fatih/structs"
	"github.com/pkg/errors"
)

type ApertureEntity struct {
	baseEntity.BaseEntity               `structs:"-"`
	FriendlyName                        string                       `json:"friendly_name" structs:"friendly_name"`
	IsFriendlyNameCustom                bool                         `json:"is_friendly_name_custom" structs:"is_friendly_name_custom" gorm:"default:false"`
	Block                               string                       `json:"block" structs:"block"`
	Facade                              string                       `json:"facade" structs:"facade" gorm:"default:''"`
	UnitName                            string                       `json:"unit_name" structs:"unit_name"`
	UnitNumber                          int                          `json:"unit_number" structs:"unit_number" gorm:"default:0"`
	Width                               int                          `json:"width" structs:"width"`
	Height                              int                          `json:"height" structs:"height"`
	Area                                float64                      `json:"area" structs:"area"`
	IsAreaCustom                        bool                         `json:"is_area_custom" structs:"is_area_custom"`
	OverallUValue                       float64                      `json:"overall_u_value" structs:"overall_u_value"`
	IsOverallUValueCustom               bool                         `json:"is_overall_u_value_custom" structs:"is_overall_u_value_custom" gorm:"default:false"`
	ActualWidth                         int64Array.Int64Array        `json:"actual_width" structs:"actual_width" gorm:"type:bigint[]"`   // One value per unit; length = Quantity
	ActualHeight                        int64Array.Int64Array        `json:"actual_height" structs:"actual_height" gorm:"type:bigint[]"` // One value per unit; length = Quantity
	Type                                *int                         `json:"type" structs:"type"`                                        // 0 = door; 1 = window
	IsTypeCustom                        bool                         `json:"is_type_custom" structs:"is_type_custom"`
	Handle                              *int                         `json:"handle" structs:"handle"` // 0 = left; 1 = right
	IsHandleCustom                      bool                         `json:"is_handle_custom" structs:"is_handle_custom"`
	ImageId                             *uint64                      `json:"image_id" structs:"image_id"`
	Image                               *fileEntity.FileEntity       `json:"-" structs:"-" gorm:"-"`
	IsImageIdCustom                     bool                         `json:"is_image_id_custom" structs:"is_image_id_custom"`
	Quantity                            int                          `json:"quantity" structs:"quantity"`
	Comment                             string                       `json:"comment" structs:"comment"`
	SupplierComment                     string                       `json:"supplier_comment" structs:"supplier_comment"`
	Accessories                         string                       `json:"accessories" structs:"accessories"`
	ColorPriceFactor                    float64                      `json:"color_price_factor" structs:"color_price_factor"`
	IsColorPriceFactorCustom            bool                         `json:"is_color_price_factor_custom" structs:"is_color_price_factor_custom"`
	GlassPriceFactor                    float64                      `json:"glass_price_factor" structs:"glass_price_factor"`
	IsGlassPriceFactorCustom            bool                         `json:"is_glass_price_factor_custom" structs:"is_glass_price_factor_custom"`
	ColorCost                           float64                      `json:"color_cost" structs:"color_cost"`
	IsColorCostCustom                   bool                         `json:"is_color_cost_custom" structs:"is_color_cost_custom"`
	GlassCost                           float64                      `json:"glass_cost" structs:"glass_cost"`
	IsGlassCostCustom                   bool                         `json:"is_glass_cost_custom" structs:"is_glass_cost_custom"`
	OtherCost                           float64                      `json:"other_cost" structs:"other_cost"`
	OtherCostComment                    string                       `json:"other_cost_comment" structs:"other_cost_comment"`
	MatrixCost                          float64                      `json:"matrix_cost" structs:"matrix_cost"` // The customised matrix cost
	IsMatrixCostCustom                  bool                         `json:"is_matrix_cost_custom" structs:"is_matrix_cost_custom"`
	NetsQuantity                        int                          `json:"nets_quantity" structs:"nets_quantity"`
	NetsPrice                           float64                      `json:"nets_price" structs:"nets_price"` // The price per net
	IsNetsPriceCustom                   bool                         `json:"is_nets_price_custom" structs:"is_nets_price_custom" gorm:"default:false"`
	NetsCost                            float64                      `json:"nets_cost" structs:"nets_cost"` // By default, nets quantity * nets price
	IsNetsCostCustom                    bool                         `json:"is_nets_cost_custom" structs:"is_nets_cost_custom" gorm:"default:false"`
	SupplierCost                        float64                      `json:"supplier_cost" structs:"supplier_cost"`
	IsSupplierCostCustom                bool                         `json:"is_supplier_cost_custom" structs:"is_supplier_cost_custom" gorm:"default:false"`
	ShippingCost                        float64                      `json:"shipping_cost" structs:"shipping_cost"`
	IsShippingCostCustom                bool                         `json:"is_shipping_cost_custom" structs:"is_shipping_cost_custom" gorm:"default:false"`
	InstallationPerUnit                 float64                      `json:"installation_per_unit" structs:"installation_per_unit" gorm:"default:0"` // per square metre
	IsInstallationPerUnitCustom         bool                         `json:"is_installation_per_unit_custom" structs:"is_installation_per_unit_custom" gorm:"default:false"`
	InstallationCost                    float64                      `json:"installation_cost" structs:"installation_cost"`
	IsInstallationCostCustom            bool                         `json:"is_installation_cost_custom" structs:"is_installation_cost_custom" gorm:"default:false"`
	ContingencyCost                     float64                      `json:"contingency_cost" structs:"contingency_cost"`
	IsContingencyCostCustom             bool                         `json:"is_contingency_cost_custom" structs:"is_contingency_cost_custom" gorm:"default:false"`
	ItemCost                            float64                      `json:"item_cost" structs:"item_cost" gorm:"default:0"` // The total cost of one item
	IsItemCostCustom                    bool                         `json:"is_item_cost_custom" structs:"is_item_cost_custom" gorm:"default:false"`
	ProfitMargin                        float64                      `json:"profit_margin" structs:"profit_margin"`
	IsProfitMarginCustom                bool                         `json:"is_profit_margin_custom" structs:"is_profit_margin_custom" gorm:"default:false"`
	SellingPrice                        float64                      `json:"selling_price" structs:"selling_price"`
	IsSellingPriceCustom                bool                         `json:"is_selling_price_custom" structs:"is_selling_price_custom" gorm:"default:false"`
	DiscountedSellingPrice              float64                      `json:"discounted_selling_price" structs:"discounted_selling_price"`
	IsDiscountedSellingPriceCustom      bool                         `json:"is_discounted_selling_price_custom" structs:"is_discounted_selling_price_custom" gorm:"default:false"`
	TotalArea                           float64                      `json:"total_area" structs:"total_area"`
	IsTotalAreaCustom                   bool                         `json:"is_total_area_custom" structs:"is_total_area_custom" gorm:"default:false"`
	TotalSellingPrice                   float64                      `json:"total_selling_price" structs:"total_selling_price"`
	IsTotalSellingPriceCustom           bool                         `json:"is_total_selling_price_custom" structs:"is_total_selling_price_custom" gorm:"default:false"`
	TotalDiscountedSellingPrice         float64                      `json:"total_discounted_selling_price" structs:"total_discounted_selling_price"`
	IsTotalDiscountedSellingPriceCustom bool                         `json:"is_total_discounted_selling_price_custom" structs:"is_total_discounted_selling_price_custom" gorm:"default:false"`
	TotalCost                           float64                      `json:"total_cost" structs:"total_cost" gorm:"default:0"` // The total cost of all items
	IsTotalCostCustom                   bool                         `json:"is_total_cost_custom" structs:"is_total_cost_custom" gorm:"default:false"`
	TotalProfit                         float64                      `json:"total_profit" structs:"total_profit" gorm:"default:0"`
	IsTotalProfitCustom                 bool                         `json:"is_total_profit_custom" structs:"is_total_profit_custom" gorm:"default:false"`
	IsCombined                          bool                         `json:"is_combined" structs:"is_combined" gorm:"-"`
	QuoteID                             *uint64                      `json:"quote_id" structs:"quote_id"`
	Quote                               *quoteEntity.QuoteEntity     `json:"-" structs:"-"`
	ProfileID                           *uint64                      `json:"profile_id" structs:"profile_id"`
	IsProfileIdCustom                   bool                         `json:"is_profile_id_custom" structs:"is_profile_id_custom"`
	Profile                             *profileEntity.ProfileEntity `json:"-" structs:"-"`
	GlassID                             *uint64                      `json:"glass_id" structs:"glass_id"`
	IsGlassIdCustom                     bool                         `json:"is_glass_id_custom" structs:"is_glass_id_custom"`
	Glass                               *glassEntity.GlassEntity     `json:"-" structs:"-"`
	ColorID                             *uint64                      `json:"color_id" structs:"color_id"`
	IsColorIdCustom                     bool                         `json:"is_color_id_custom" structs:"is_color_id_custom"`
	Color                               *colorEntity.ColorEntity     `json:"-" structs:"-"`
	ParentApertureID                    *uint64                      `json:"parent_aperture_id" structs:"parent_aperture_id"`
	ParentAperture                      *ApertureEntity              `json:"-" structs:"-"`
	CreateUser                          *userEntity.UserEntity       `json:"-" structs:"-" gorm:"foreignKey:created_by"` // The user who created the aperture
	UpdateUser                          *userEntity.UserEntity       `json:"-" structs:"-" gorm:"foreignKey:updated_by"` // The user who last updated the aperture
	ChildApertures                      []*ApertureEntity            `json:"-" structs:"-" gorm:"-"`
}

func New() *ApertureEntity {
	entity := new(ApertureEntity)
	entity.init()
	return entity
}

func (e *ApertureEntity) init() {
	// Intentionally left empty
}

func (e *ApertureEntity) GetFriendlyName() string {
	return e.FriendlyName
}
func (e *ApertureEntity) GetIsFriendlyNameCustom() bool {
	return e.IsFriendlyNameCustom
}
func (e *ApertureEntity) GetBlock() string {
	return e.Block
}
func (e *ApertureEntity) GetFacade() string {
	return e.Facade
}
func (e *ApertureEntity) GetUnitName() string {
	return e.UnitName
}
func (e *ApertureEntity) GetUnitNumber() int {
	return e.UnitNumber
}
func (e *ApertureEntity) GetUnitFullName() string {
	if e.UnitNumber == 0 {
		return e.UnitName
	}

	return fmt.Sprintf("%s %d", e.UnitName, e.UnitNumber)
}
func (e *ApertureEntity) GetWidth() int {
	return e.Width
}
func (e *ApertureEntity) GetHeight() int {
	return e.Height
}
func (e *ApertureEntity) GetArea() float64 {
	return e.Area
}
func (e *ApertureEntity) GetIsAreaCustom() bool {
	return e.IsAreaCustom
}
func (e *ApertureEntity) GetOverallUValue() float64 {
	return e.OverallUValue
}
func (e *ApertureEntity) GetIsOverallUValueCustom() bool {
	return e.IsOverallUValueCustom
}
func (e *ApertureEntity) GetActualWidth() int64Array.Int64Array {
	return e.ActualWidth
}
func (e *ApertureEntity) GetActualHeight() int64Array.Int64Array {
	return e.ActualHeight
}
func (e *ApertureEntity) GetType() *int {
	return e.Type
}
func (e *ApertureEntity) GetTypeName() *apertureType.ApertureType {
	key := e.GetType()
	if key == nil {
		return nil
	}
	return apertureType.FromInt(*key)
}
func (e *ApertureEntity) GetIsTypeCustom() bool {
	return e.IsTypeCustom
}
func (e *ApertureEntity) GetHandle() *int {
	return e.Handle
}
func (e *ApertureEntity) GetHandleName() *handleType.HandleType {
	key := e.GetHandle()
	if key == nil {
		return nil
	}
	return handleType.FromInt(*key)
}
func (e *ApertureEntity) GetIsHandleCustom() bool {
	return e.IsHandleCustom
}
func (e *ApertureEntity) GetImageId() *uint64 {
	return e.ImageId
}
func (e *ApertureEntity) GetIsImageIdCustom() bool {
	return e.IsImageIdCustom
}
func (e *ApertureEntity) GetQuantity() int {
	return e.Quantity
}
func (e *ApertureEntity) GetComment() string {
	return e.Comment
}
func (e *ApertureEntity) GetSupplierComment() string {
	return e.SupplierComment
}
func (e *ApertureEntity) GetAccessories() string {
	return e.Accessories
}
func (e *ApertureEntity) GetColorPriceFactor() float64 {
	return e.ColorPriceFactor
}
func (e *ApertureEntity) GetIsColorPriceFactorCustom() bool {
	return e.IsColorPriceFactorCustom
}
func (e *ApertureEntity) GetGlassPriceFactor() float64 {
	return e.GlassPriceFactor
}
func (e *ApertureEntity) GetIsGlassPriceFactorCustom() bool {
	return e.IsGlassPriceFactorCustom
}
func (e *ApertureEntity) GetColorCost() float64 {
	return e.ColorCost
}
func (e *ApertureEntity) GetIsColorCostCustom() bool {
	return e.IsColorCostCustom
}
func (e *ApertureEntity) GetGlassCost() float64 {
	return e.GlassCost
}
func (e *ApertureEntity) GetIsGlassCostCustom() bool {
	return e.IsGlassCostCustom
}
func (e *ApertureEntity) GetOtherCost() float64 {
	return e.OtherCost
}
func (e *ApertureEntity) GetOtherCostComment() string {
	return e.OtherCostComment
}
func (e *ApertureEntity) GetMatrixCost() float64 {
	return e.MatrixCost
}
func (e *ApertureEntity) GetIsMatrixCostCustom() bool {
	return e.IsMatrixCostCustom
}
func (e *ApertureEntity) GeNetsQuantity() int {
	return e.NetsQuantity
}
func (e *ApertureEntity) GetNetsPrice() float64 {
	return e.NetsPrice
}
func (e *ApertureEntity) GetIsNetsPriceCustom() bool {
	return e.IsNetsPriceCustom
}
func (e *ApertureEntity) GetNetsCost() float64 {
	return e.NetsCost
}
func (e *ApertureEntity) GetIsNetsCostCustom() bool {
	return e.IsNetsCostCustom
}
func (e *ApertureEntity) GetSupplierCost() float64 {
	return e.SupplierCost
}
func (e *ApertureEntity) GetIsSupplierCostCustom() bool {
	return e.IsSupplierCostCustom
}
func (e *ApertureEntity) GetShippingCost() float64 {
	return e.ShippingCost
}
func (e *ApertureEntity) GetIsShippingCostCustom() bool {
	return e.IsShippingCostCustom
}
func (e *ApertureEntity) GetInstallationPerUnit() float64 {
	return e.InstallationPerUnit
}
func (e *ApertureEntity) GetIsInstallationPerUnitCustom() bool {
	return e.IsInstallationPerUnitCustom
}
func (e *ApertureEntity) GetInstallationCost() float64 {
	return e.InstallationCost
}
func (e *ApertureEntity) GetIsInstallationCostCustom() bool {
	return e.IsInstallationCostCustom
}
func (e *ApertureEntity) GetContingencyCost() float64 {
	return e.ContingencyCost
}
func (e *ApertureEntity) GetIsContingencyCostCustom() bool {
	return e.IsContingencyCostCustom
}
func (e *ApertureEntity) GetItemCost() float64 {
	return e.ItemCost
}
func (e *ApertureEntity) GetIsItemCostCustom() bool {
	return e.IsItemCostCustom
}
func (e *ApertureEntity) GetProfitMargin() float64 {
	return e.ProfitMargin
}
func (e *ApertureEntity) GetIsProfitMarginCustom() bool {
	return e.IsProfitMarginCustom
}
func (e *ApertureEntity) GetSellingPrice() float64 {
	return e.SellingPrice
}
func (e *ApertureEntity) GetIsSellingPriceCustom() bool {
	return e.IsSellingPriceCustom
}
func (e *ApertureEntity) GetDiscountedSellingPrice() float64 {
	return e.DiscountedSellingPrice
}
func (e *ApertureEntity) GetIsDiscountedSellingPriceCustom() bool {
	return e.IsDiscountedSellingPriceCustom
}
func (e *ApertureEntity) GetTotalArea() float64 {
	return e.TotalArea
}
func (e *ApertureEntity) GetIsTotalAreaCustom() bool {
	return e.IsTotalAreaCustom
}
func (e *ApertureEntity) GetTotalSellingPrice() float64 {
	return e.TotalSellingPrice
}
func (e *ApertureEntity) GetIsTotalSellingPriceCustom() bool {
	return e.IsTotalSellingPriceCustom
}
func (e *ApertureEntity) GetTotalDiscountedSellingPrice() float64 {
	return e.TotalDiscountedSellingPrice
}
func (e *ApertureEntity) GetIsTotalDiscountedSellingPriceCustom() bool {
	return e.IsTotalDiscountedSellingPriceCustom
}
func (e *ApertureEntity) GetTotalCost() float64 {
	return e.TotalCost
}
func (e *ApertureEntity) GetIsTotalCostCustom() bool {
	return e.IsTotalCostCustom
}
func (e *ApertureEntity) GetTotalProfit() float64 {
	return e.TotalProfit
}
func (e *ApertureEntity) GetIsTotalProfitCustom() bool {
	return e.IsTotalProfitCustom
}
func (e *ApertureEntity) GetIsCombined() bool {
	return e.IsCombined
}
func (e *ApertureEntity) GetQuoteID() *uint64 {
	return e.QuoteID
}
func (e *ApertureEntity) GetProfileID() *uint64 {
	return e.ProfileID
}
func (e *ApertureEntity) GetIsProfileIDCustom() bool {
	return e.IsProfileIdCustom
}
func (e *ApertureEntity) GetGlassID() *uint64 {
	return e.GlassID
}
func (e *ApertureEntity) GetIsGlassIDCustom() bool {
	return e.IsGlassIdCustom
}
func (e *ApertureEntity) GetColorID() *uint64 {
	return e.ColorID
}
func (e *ApertureEntity) GetIsColorIDCustom() bool {
	return e.IsColorIdCustom
}
func (e *ApertureEntity) GetParentApertureID() *uint64 {
	return e.ParentApertureID
}

func (e *ApertureEntity) SetFriendlyName(value string) {
	e.FriendlyName = value
}
func (e *ApertureEntity) SetIsFriendlyNameCustom(value bool) {
	e.IsFriendlyNameCustom = value
}
func (e *ApertureEntity) SetBlock(value string) {
	e.Block = value
}
func (e *ApertureEntity) SetFacade(value string) {
	e.Facade = value
}
func (e *ApertureEntity) SetUnitName(value string) {
	e.UnitName = value
}
func (e *ApertureEntity) SetUnitNumber(value int) {
	e.UnitNumber = value
}
func (e *ApertureEntity) SetWidth(value int) {
	e.Width = value
}
func (e *ApertureEntity) SetHeight(value int) {
	e.Height = value
}
func (e *ApertureEntity) SetArea(value float64) {
	e.Area = value
}
func (e *ApertureEntity) SetIsAreaCustom(value bool) {
	e.IsAreaCustom = value
}
func (e *ApertureEntity) SetOverallUValue(value float64) {
	e.OverallUValue = value
}
func (e *ApertureEntity) SetIsOverallUValueCustom(value bool) {
	e.IsOverallUValueCustom = value
}
func (e *ApertureEntity) SetActualWidth(value int64Array.Int64Array) {
	e.ActualWidth = value
}
func (e *ApertureEntity) SetActualHeight(value int64Array.Int64Array) {
	e.ActualHeight = value
}
func (e *ApertureEntity) SetType(value *int) {
	e.Type = value
}
func (e *ApertureEntity) SetTypeName(value *apertureType.ApertureType) {
	if value == nil {
		e.SetType(nil)
	} else {
		t := int(*value)
		e.SetType(&t)
	}
}
func (e *ApertureEntity) SetIsTypeCustom(value bool) {
	e.IsTypeCustom = value
}
func (e *ApertureEntity) SetHandle(value *int) {
	e.Handle = value
}
func (e *ApertureEntity) SetHandleName(value *handleType.HandleType) {
	if value == nil {
		e.SetHandle(nil)
	} else {
		t := int(*value)
		e.SetHandle(&t)
	}
}
func (e *ApertureEntity) SetIsHandleCustom(value bool) {
	e.IsHandleCustom = value
}
func (e *ApertureEntity) SetImageId(value *uint64) {
	e.ImageId = value
}
func (e *ApertureEntity) SetIsImageIdCustom(value bool) {
	e.IsImageIdCustom = value
}
func (e *ApertureEntity) SetQuantity(value int) {
	e.Quantity = value
}
func (e *ApertureEntity) SetComment(value string) {
	e.Comment = value
}
func (e *ApertureEntity) SetSupplierComment(value string) {
	e.SupplierComment = value
}
func (e *ApertureEntity) SetAccessories(value string) {
	e.Accessories = value
}
func (e *ApertureEntity) SetColorPriceFactor(value float64) {
	e.ColorPriceFactor = value
}
func (e *ApertureEntity) SetIsColorPriceFactorCustom(value bool) {
	e.IsColorPriceFactorCustom = value
}
func (e *ApertureEntity) SetGlassPriceFactor(value float64) {
	e.GlassPriceFactor = value
}
func (e *ApertureEntity) SetIsGlassPriceFactorCustom(value bool) {
	e.IsGlassPriceFactorCustom = value
}
func (e *ApertureEntity) SetColorCost(value float64) {
	e.ColorCost = value
}
func (e *ApertureEntity) SetIsColorCostCustom(value bool) {
	e.IsColorCostCustom = value
}
func (e *ApertureEntity) SetGlassCost(value float64) {
	e.GlassCost = value
}
func (e *ApertureEntity) SetIsGlassCostCustom(value bool) {
	e.IsGlassCostCustom = value
}
func (e *ApertureEntity) SetOtherCost(value float64) {
	e.OtherCost = value
}
func (e *ApertureEntity) SetOtherCostComment(value string) {
	e.OtherCostComment = value
}
func (e *ApertureEntity) SetMatrixCost(value float64) {
	e.MatrixCost = value
}
func (e *ApertureEntity) SetIsMatrixCostCustom(value bool) {
	e.IsMatrixCostCustom = value
}
func (e *ApertureEntity) SetNetsQuantity(value int) {
	e.NetsQuantity = value
}
func (e *ApertureEntity) SetNetsPrice(value float64) {
	e.NetsPrice = value
}
func (e *ApertureEntity) SetIsNetsPriceCustom(value bool) {
	e.IsNetsPriceCustom = value
}
func (e *ApertureEntity) SetNetsCost(value float64) {
	e.NetsCost = value
}
func (e *ApertureEntity) SetIsNetsCostCustom(value bool) {
	e.IsNetsCostCustom = value
}
func (e *ApertureEntity) SetSupplierCost(value float64) {
	e.SupplierCost = value
}
func (e *ApertureEntity) SetIsSupplierCostCustom(value bool) {
	e.IsSupplierCostCustom = value
}
func (e *ApertureEntity) SetShippingCost(value float64) {
	e.ShippingCost = value
}
func (e *ApertureEntity) SetIsShippingCostCustom(value bool) {
	e.IsShippingCostCustom = value
}
func (e *ApertureEntity) SetInstallationPerUnit(value float64) {
	e.InstallationPerUnit = value
}
func (e *ApertureEntity) SetIsInstallationPerUnitCustom(value bool) {
	e.IsInstallationPerUnitCustom = value
}
func (e *ApertureEntity) SetInstallationCost(value float64) {
	e.InstallationCost = value
}
func (e *ApertureEntity) SetIsInstallationCostCustom(value bool) {
	e.IsInstallationCostCustom = value
}
func (e *ApertureEntity) SetContingencyCost(value float64) {
	e.ContingencyCost = value
}
func (e *ApertureEntity) SetIsContingencyCostCustom(value bool) {
	e.IsContingencyCostCustom = value
}
func (e *ApertureEntity) SetItemCost(value float64) {
	e.ItemCost = value
}
func (e *ApertureEntity) SetIsItemCostCustom(value bool) {
	e.IsItemCostCustom = value
}
func (e *ApertureEntity) SetProfitMargin(value float64) {
	e.ProfitMargin = value
}
func (e *ApertureEntity) SetIsProfitMarginCustom(value bool) {
	e.IsProfitMarginCustom = value
}
func (e *ApertureEntity) SetSellingPrice(value float64) {
	e.SellingPrice = value
}
func (e *ApertureEntity) SetIsSellingPriceCustom(value bool) {
	e.IsSellingPriceCustom = value
}
func (e *ApertureEntity) SetDiscountedSellingPrice(value float64) {
	e.DiscountedSellingPrice = value
}
func (e *ApertureEntity) SetIsDiscountedSellingPriceCustom(value bool) {
	e.IsDiscountedSellingPriceCustom = value
}
func (e *ApertureEntity) SetTotalArea(value float64) {
	e.TotalArea = value
}
func (e *ApertureEntity) SetIsTotalAreaCustom(value bool) {
	e.IsTotalAreaCustom = value
}
func (e *ApertureEntity) SetTotalSellingPrice(value float64) {
	e.TotalSellingPrice = value
}
func (e *ApertureEntity) SetIsTotalSellingPriceCustom(value bool) {
	e.IsTotalSellingPriceCustom = value
}
func (e *ApertureEntity) SetTotalDiscountedSellingPrice(value float64) {
	e.TotalDiscountedSellingPrice = value
}
func (e *ApertureEntity) SetIsTotalDiscountedSellingPriceCustom(value bool) {
	e.IsTotalDiscountedSellingPriceCustom = value
}
func (e *ApertureEntity) SetTotalCost(value float64) {
	e.TotalCost = value
}
func (e *ApertureEntity) SetIsTotalCostCustom(value bool) {
	e.IsTotalCostCustom = value
}
func (e *ApertureEntity) SetTotalProfit(value float64) {
	e.TotalProfit = value
}
func (e *ApertureEntity) SetIsTotalProfitCustom(value bool) {
	e.IsTotalProfitCustom = value
}
func (e *ApertureEntity) SetQuoteID(value *uint64) {
	e.QuoteID = value
}
func (e *ApertureEntity) SetProfileID(value *uint64) {
	e.ProfileID = value
}
func (e *ApertureEntity) SetIsProfileIDCustom(value bool) {
	e.IsProfileIdCustom = value
}
func (e *ApertureEntity) SetGlassID(value *uint64) {
	e.GlassID = value
}
func (e *ApertureEntity) SetIsGlassIDCustom(value bool) {
	e.IsGlassIdCustom = value
}
func (e *ApertureEntity) SetColorID(value *uint64) {
	e.ColorID = value
}
func (e *ApertureEntity) SetIsColorIDCustom(value bool) {
	e.IsColorIdCustom = value
}
func (e *ApertureEntity) SetParentApertureID(value *uint64) {
	e.ParentApertureID = value
}

func (*ApertureEntity) TableName() string {
	return TableName()
}
func TableName() string {
	return "apertures"
}

func (e *ApertureEntity) FieldNames() []string {
	result := e.BaseEntity.FieldNames()
	return append(result, []string{"is_friendly_name_custom", "block", "facade", "unit_name", "unit_number", "width", "height", "is_area_custom", "is_overall_u_value_custom", "actual_width", "actual_height", "is_type_custom", "is_handle_custom", "is_image_id_custom", "quantity", "comment", "supplier_comment", "accessories", "is_color_price_factor_custom", "is_glass_price_factor_custom", "is_color_cost_custom", "is_glass_cost_custom", "is_matrix_cost_custom", "is_nets_price_custom", "is_nets_cost_custom", "is_supplier_cost_custom", "is_shipping_cost_custom", "is_installation_per_unit_custom", "is_installation_cost_custom", "is_contingency_cost_custom", "is_item_cost_custom", "is_profit_margin_custom", "is_selling_price_custom", "is_discounted_selling_price_custom", "is_total_area_custom", "is_total_selling_price_custom", "is_total_discounted_selling_price_custom", "is_total_cost_custom", "is_total_profit_custom", "quote_id", "is_profile_id_custom", "is_glass_id_custom", "is_color_id_custom", "parent_aperture_id"}...)
}

func (e *ApertureEntity) ToJson() map[string]any {
	json := structs.Map(e)

	for key, element := range e.BaseEntity.ToJson() {
		json[key] = element
	}

	// Always return arrays (never null) for the per-unit dimensions
	json["actual_width"] = nonNilArray(e.ActualWidth)
	json["actual_height"] = nonNilArray(e.ActualHeight)

	if e.Quote != nil {
		json["quote"] = e.Quote.ToJson()
	}
	if e.Profile != nil {
		json["profile"] = e.Profile.ToJson()
	}
	if e.Glass != nil {
		json["glass"] = e.Glass.ToJson()
	}
	if e.Color != nil {
		json["color"] = e.Color.ToJson()
	}
	if e.ParentAperture != nil {
		json["parent_aperture"] = e.ParentAperture.ToJson()
	}

	if e.CreateUser != nil {
		json["create_user"] = e.CreateUser.ToJson()
	}
	if e.UpdateUser != nil {
		json["update_user"] = e.UpdateUser.ToJson()
	}

	return json
}

func nonNilArray(value int64Array.Int64Array) int64Array.Int64Array {
	if value == nil {
		return int64Array.Int64Array{}
	}
	return value
}

// GetActualAreaSqm returns the summed area (square metres) of all units, using each unit's own actual width and height.
func (e *ApertureEntity) GetActualAreaSqm() float64 {
	return int64Array.AreaSqm(e.ActualWidth, e.ActualHeight)
}

func (e *ApertureEntity) FillFromJson(json map[string]any) error {
	if err := e.BaseEntity.FillFromJson(json); err != nil {
		return err
	}

	_, err := mapUtils.FieldAnyToType(json, "unit_name", fmt.Sprintf("%T", e.UnitName))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "unit_number", fmt.Sprintf("%T", e.UnitNumber))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "width", fmt.Sprintf("%T", e.Width))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "height", fmt.Sprintf("%T", e.Height))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "area", fmt.Sprintf("%T", e.Area))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "overall_u_value", fmt.Sprintf("%T", e.OverallUValue))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "actual_width", fmt.Sprintf("%T", e.ActualWidth))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "actual_height", fmt.Sprintf("%T", e.ActualHeight))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "type", fmt.Sprintf("%T", e.Type))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "handle", fmt.Sprintf("%T", e.Handle))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "image_id", fmt.Sprintf("%T", e.ImageId))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "quantity", fmt.Sprintf("%T", e.Quantity))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "color_price_factor", fmt.Sprintf("%T", e.ColorPriceFactor))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "glass_price_factor", fmt.Sprintf("%T", e.GlassPriceFactor))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "color_cost", fmt.Sprintf("%T", e.ColorCost))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "glass_cost", fmt.Sprintf("%T", e.GlassCost))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "other_cost", fmt.Sprintf("%T", e.OtherCost))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "other_cost_comment", fmt.Sprintf("%T", e.OtherCostComment))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "matrix_cost", fmt.Sprintf("%T", e.MatrixCost))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "nets_quantity", fmt.Sprintf("%T", e.NetsQuantity))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "nets_price", fmt.Sprintf("%T", e.NetsPrice))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "nets_cost", fmt.Sprintf("%T", e.NetsCost))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "supplier_cost", fmt.Sprintf("%T", e.SupplierCost))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "shipping_cost", fmt.Sprintf("%T", e.ShippingCost))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "installation_per_unit", fmt.Sprintf("%T", e.InstallationPerUnit))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "installation_cost", fmt.Sprintf("%T", e.InstallationCost))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "contingency_cost", fmt.Sprintf("%T", e.ContingencyCost))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "item_cost", fmt.Sprintf("%T", e.ItemCost))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "profit_margin", fmt.Sprintf("%T", e.ProfitMargin))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "selling_price", fmt.Sprintf("%T", e.SellingPrice))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "discounted_selling_price", fmt.Sprintf("%T", e.DiscountedSellingPrice))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "total_area", fmt.Sprintf("%T", e.TotalArea))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "total_selling_price", fmt.Sprintf("%T", e.TotalSellingPrice))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "total_discounted_selling_price", fmt.Sprintf("%T", e.TotalDiscountedSellingPrice))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "total_cost", fmt.Sprintf("%T", e.TotalCost))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "total_profit", fmt.Sprintf("%T", e.TotalProfit))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "quote_id", fmt.Sprintf("%T", e.QuoteID))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "profile_id", fmt.Sprintf("%T", e.ProfileID))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "glass_id", fmt.Sprintf("%T", e.GlassID))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "color_id", fmt.Sprintf("%T", e.ColorID))
	if err != nil {
		return err
	}
	_, err = mapUtils.FieldAnyToType(json, "parent_aperture_id", fmt.Sprintf("%T", e.ParentApertureID))
	if err != nil {
		return err
	}

	err = mapUtils.ToStruct_usingJsonTags(json, e)
	if err != nil {
		return err
	}

	if mapUtils.HasKey(json, "color") {
		color := colorEntity.New()
		color.FillFromJson(json["color"].(map[string]any))
		e.Color = color
	}
	if mapUtils.HasKey(json, "create_user") {
		user := userEntity.New()
		user.FillFromJson(json["create_user"].(map[string]any))
		e.CreateUser = user
	}
	if mapUtils.HasKey(json, "glass") {
		glass := glassEntity.New()
		glass.FillFromJson(json["glass"].(map[string]any))
		e.Glass = glass
	}
	if mapUtils.HasKey(json, "profile") {
		profile := profileEntity.New()
		profile.FillFromJson(json["profile"].(map[string]any))
		e.Profile = profile
	}
	if mapUtils.HasKey(json, "update_user") {
		user := userEntity.New()
		user.FillFromJson(json["update_user"].(map[string]any))
		e.UpdateUser = user
	}

	// Make sure Accesssories has no blank lines
	{
		lines := strings.Split(e.GetAccessories(), "\n")

		accessories := ""
		for _, line := range lines {
			if line != "" {
				accessories += line + "\n"
			}
		}
		accessories = strings.TrimRight(accessories, "\n")

		e.SetAccessories(accessories)
	}

	return nil
}

func (e *ApertureEntity) FillFromDb(dbMap map[string]any) error {
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

	if e.QuoteID != nil {
		// This Aperture has a quote. Try to fill in the quote's details
		itemMap := map[string]any{}
		quote := quoteEntity.New()
		keyPrefix := quote.TableName() + "_"
		for key, value := range dbMap {
			if strings.HasPrefix(key, keyPrefix) {
				newKey := strings.TrimPrefix(key, keyPrefix)
				itemMap[newKey] = value
			}
		}

		if len(itemMap) > 0 {
			if err := quote.FillFromJson(itemMap); err != nil {
				return err
			}

			e.Quote = quote
		}
	}

	if e.ProfileID != nil {
		// This Aperture has a profile. Try to fill in the profile's details
		itemMap := map[string]any{}
		profile := profileEntity.New()
		keyPrefix := profile.TableName() + "_"
		for key, value := range dbMap {
			if strings.HasPrefix(key, keyPrefix) {
				newKey := strings.TrimPrefix(key, keyPrefix)
				itemMap[newKey] = value
			}
		}

		if len(itemMap) > 0 {
			if err := profile.FillFromJson(itemMap); err != nil {
				return err
			}

			e.Profile = profile
		}
	}

	if e.GlassID != nil {
		// This Aperture has a glass. Try to fill in the glass' details
		itemMap := map[string]any{}
		glass := glassEntity.New()
		keyPrefix := glass.TableName() + "_"
		for key, value := range dbMap {
			if strings.HasPrefix(key, keyPrefix) {
				newKey := strings.TrimPrefix(key, keyPrefix)
				itemMap[newKey] = value
			}
		}

		if len(itemMap) > 0 {
			if err := glass.FillFromJson(itemMap); err != nil {
				return err
			}

			e.Glass = glass
		}
	}

	if e.ColorID != nil {
		// This Aperture has a color. Try to fill in the color's details
		itemMap := map[string]any{}
		color := colorEntity.New()
		keyPrefix := color.TableName() + "_"
		for key, value := range dbMap {
			if strings.HasPrefix(key, keyPrefix) {
				newKey := strings.TrimPrefix(key, keyPrefix)
				itemMap[newKey] = value
			}
		}

		if len(itemMap) > 0 {
			if err := color.FillFromJson(itemMap); err != nil {
				return err
			}

			e.Color = color
		}
	}

	if e.ParentApertureID != nil {
		// This Profile has a supplier. Try to fill in the supplier's details
		itemMap := map[string]any{}
		parentAperture := New()
		keyPrefix := parentAperture.TableName() + "_"
		for key, value := range dbMap {
			if strings.HasPrefix(key, keyPrefix) {
				newKey := strings.TrimPrefix(key, keyPrefix)
				itemMap[newKey] = value
			}
		}

		if len(itemMap) > 0 {
			if err := parentAperture.FillFromJson(itemMap); err != nil {
				return err
			}

			e.ParentAperture = parentAperture
		}
	}

	{
		// This Aperture has a create user. Try to fill in the user's details
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
		// This Aperture has an update user. Try to fill in the user's details
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

func (e *ApertureEntity) Clone() (baseEntity.BaseEntityI, error) {
	clone := New()
	err := clone.FillFromJson(e.ToJson())
	return clone, err
}

func (e *ApertureEntity) Validate() error {
	// Validate the Type
	{
		key := e.GetType()
		if key != nil && !apertureType.Exists(*key) {
			return errors.Errorf("Type '%d' is invalid", *key)
		}
	}

	// Validate the Handle
	{
		key := e.GetType()
		if key != nil && !handleType.Exists(*key) {
			return errors.Errorf("Handle '%d' is invalid", *key)
		}
	}

	// Validate the actual dimensions: one value per unit
	if len(e.ActualWidth) != e.Quantity {
		return errors.Errorf("actual_width must contain exactly %d value(s) (one per unit), but has %d", e.Quantity, len(e.ActualWidth))
	}
	if len(e.ActualHeight) != e.Quantity {
		return errors.Errorf("actual_height must contain exactly %d value(s) (one per unit), but has %d", e.Quantity, len(e.ActualHeight))
	}

	// Validate the Accessories
	{
		lines := strings.Split(e.GetAccessories(), "\n")
		if len(lines) > 3 {
			return errors.Errorf("You cannot submit more than 3 accessories")
		}
	}
	return nil
}

func (e *ApertureEntity) GetUpdates(oldItem baseEntity.BaseEntityI) map[string]any {
	oldAperture := oldItem.(*ApertureEntity)
	updates := map[string]any{}

	if e.FriendlyName != oldAperture.FriendlyName {
		updates["friendly_name"] = e.FriendlyName
	}
	if e.IsFriendlyNameCustom != oldAperture.IsFriendlyNameCustom {
		updates["is_friendly_name_custom"] = e.IsFriendlyNameCustom
	}
	if e.Block != oldAperture.Block {
		updates["block"] = e.Block
	}
	if e.Facade != oldAperture.Facade {
		updates["facade"] = e.Facade
	}
	if e.UnitName != oldAperture.UnitName {
		updates["unit_name"] = e.UnitName
	}
	if e.UnitNumber != oldAperture.UnitNumber {
		updates["unit_number"] = e.UnitNumber
	}
	if e.Width != oldAperture.Width {
		updates["width"] = e.Width
	}
	if e.Height != oldAperture.Height {
		updates["height"] = e.Height
	}
	if e.Area != oldAperture.Area {
		updates["area"] = e.Area
	}
	if e.IsAreaCustom != oldAperture.IsAreaCustom {
		updates["is_area_custom"] = e.IsAreaCustom
	}
	if e.OverallUValue != oldAperture.OverallUValue {
		updates["overall_u_value"] = e.OverallUValue
	}
	if e.IsOverallUValueCustom != oldAperture.IsOverallUValueCustom {
		updates["is_overall_u_value_custom"] = e.IsOverallUValueCustom
	}
	if !slices.Equal(e.ActualWidth, oldAperture.ActualWidth) {
		updates["actual_width"] = e.ActualWidth
	}
	if !slices.Equal(e.ActualHeight, oldAperture.ActualHeight) {
		updates["actual_height"] = e.ActualHeight
	}
	if e.Type != oldAperture.Type {
		updates["type"] = e.Type
	}
	if e.IsTypeCustom != oldAperture.IsTypeCustom {
		updates["is_type_custom"] = e.IsTypeCustom
	}
	if e.Handle != oldAperture.Handle {
		updates["handle"] = e.Handle
	}
	if e.IsHandleCustom != oldAperture.IsHandleCustom {
		updates["is_handle_custom"] = e.IsHandleCustom
	}
	if e.ImageId != oldAperture.ImageId {
		updates["image_id"] = e.ImageId
	}
	if e.IsImageIdCustom != oldAperture.IsImageIdCustom {
		updates["is_image_id_custom"] = e.IsImageIdCustom
	}
	if e.Quantity != oldAperture.Quantity {
		updates["quantity"] = e.Quantity
	}
	if e.Comment != oldAperture.Comment {
		updates["comment"] = e.Comment
	}
	if e.SupplierComment != oldAperture.SupplierComment {
		updates["supplier_comment"] = e.SupplierComment
	}
	if e.Accessories != oldAperture.Accessories {
		updates["accessories"] = e.Accessories
	}
	if e.ColorPriceFactor != oldAperture.ColorPriceFactor {
		updates["color_price_factor"] = e.ColorPriceFactor
	}
	if e.IsColorPriceFactorCustom != oldAperture.IsColorPriceFactorCustom {
		updates["is_color_price_factor_custom"] = e.IsColorPriceFactorCustom
	}
	if e.GlassPriceFactor != oldAperture.GlassPriceFactor {
		updates["glass_price_factor"] = e.GlassPriceFactor
	}
	if e.IsGlassPriceFactorCustom != oldAperture.IsGlassPriceFactorCustom {
		updates["is_glass_price_factor_custom"] = e.IsGlassPriceFactorCustom
	}
	if e.ColorCost != oldAperture.ColorCost {
		updates["color_cost"] = e.ColorCost
	}
	if e.IsColorCostCustom != oldAperture.IsColorCostCustom {
		updates["is_color_cost_custom"] = e.IsColorCostCustom
	}
	if e.GlassCost != oldAperture.GlassCost {
		updates["glass_cost"] = e.GlassCost
	}
	if e.IsGlassCostCustom != oldAperture.IsGlassCostCustom {
		updates["is_glass_cost_custom"] = e.IsGlassCostCustom
	}
	if e.OtherCost != oldAperture.OtherCost {
		updates["other_cost"] = e.OtherCost
	}
	if e.OtherCostComment != oldAperture.OtherCostComment {
		updates["other_cost_comment"] = e.OtherCostComment
	}
	if e.MatrixCost != oldAperture.MatrixCost {
		updates["matrix_cost"] = e.MatrixCost
	}
	if e.IsMatrixCostCustom != oldAperture.IsMatrixCostCustom {
		updates["is_matrix_cost_custom"] = e.IsMatrixCostCustom
	}
	if e.NetsQuantity != oldAperture.NetsQuantity {
		updates["nets_quantity"] = e.NetsQuantity
	}
	if e.NetsPrice != oldAperture.NetsPrice {
		updates["nets_price"] = e.NetsPrice
	}
	if e.IsNetsPriceCustom != oldAperture.IsNetsPriceCustom {
		updates["is_nets_price_custom"] = e.IsNetsPriceCustom
	}
	if e.NetsCost != oldAperture.NetsCost {
		updates["nets_cost"] = e.NetsCost
	}
	if e.IsNetsCostCustom != oldAperture.IsNetsCostCustom {
		updates["is_nets_cost_custom"] = e.IsNetsCostCustom
	}
	if e.SupplierCost != oldAperture.SupplierCost {
		updates["supplier_cost"] = e.SupplierCost
	}
	if e.IsSupplierCostCustom != oldAperture.IsSupplierCostCustom {
		updates["is_supplier_cost_custom"] = e.IsSupplierCostCustom
	}
	if e.ShippingCost != oldAperture.ShippingCost {
		updates["shipping_cost"] = e.ShippingCost
	}
	if e.IsShippingCostCustom != oldAperture.IsShippingCostCustom {
		updates["is_shipping_cost_custom"] = e.IsShippingCostCustom
	}
	if e.InstallationPerUnit != oldAperture.InstallationPerUnit {
		updates["installation_per_unit"] = e.InstallationPerUnit
	}
	if e.IsInstallationPerUnitCustom != oldAperture.IsInstallationPerUnitCustom {
		updates["is_installation_per_unit_custom"] = e.IsInstallationPerUnitCustom
	}
	if e.InstallationCost != oldAperture.InstallationCost {
		updates["installation_cost"] = e.InstallationCost
	}
	if e.IsInstallationCostCustom != oldAperture.IsInstallationCostCustom {
		updates["is_installation_cost_custom"] = e.IsInstallationCostCustom
	}
	if e.ContingencyCost != oldAperture.ContingencyCost {
		updates["contingency_cost"] = e.ContingencyCost
	}
	if e.IsContingencyCostCustom != oldAperture.IsContingencyCostCustom {
		updates["is_contingency_cost_custom"] = e.IsContingencyCostCustom
	}
	if e.ItemCost != oldAperture.ItemCost {
		updates["item_cost"] = e.ItemCost
	}
	if e.IsItemCostCustom != oldAperture.IsItemCostCustom {
		updates["is_item_cost_custom"] = e.IsItemCostCustom
	}
	if e.ProfitMargin != oldAperture.ProfitMargin {
		updates["profit_margin"] = e.ProfitMargin
	}
	if e.IsProfitMarginCustom != oldAperture.IsProfitMarginCustom {
		updates["is_profit_margin_custom"] = e.IsProfitMarginCustom
	}
	if e.SellingPrice != oldAperture.SellingPrice {
		updates["selling_price"] = e.SellingPrice
	}
	if e.IsSellingPriceCustom != oldAperture.IsSellingPriceCustom {
		updates["is_selling_price_custom"] = e.IsSellingPriceCustom
	}
	if e.DiscountedSellingPrice != oldAperture.DiscountedSellingPrice {
		updates["discounted_selling_price"] = e.DiscountedSellingPrice
	}
	if e.IsDiscountedSellingPriceCustom != oldAperture.IsDiscountedSellingPriceCustom {
		updates["is_discounted_selling_price_custom"] = e.IsDiscountedSellingPriceCustom
	}
	if e.TotalArea != oldAperture.TotalArea {
		updates["total_area"] = e.TotalArea
	}
	if e.IsTotalAreaCustom != oldAperture.IsTotalAreaCustom {
		updates["is_total_area_custom"] = e.IsTotalAreaCustom
	}
	if e.TotalSellingPrice != oldAperture.TotalSellingPrice {
		updates["total_selling_price"] = e.TotalSellingPrice
	}
	if e.IsTotalSellingPriceCustom != oldAperture.IsTotalSellingPriceCustom {
		updates["is_total_selling_price_custom"] = e.IsTotalSellingPriceCustom
	}
	if e.TotalDiscountedSellingPrice != oldAperture.TotalDiscountedSellingPrice {
		updates["total_discounted_selling_price"] = e.TotalDiscountedSellingPrice
	}
	if e.IsTotalDiscountedSellingPriceCustom != oldAperture.IsTotalDiscountedSellingPriceCustom {
		updates["is_total_discounted_selling_price_custom"] = e.IsTotalDiscountedSellingPriceCustom
	}
	if e.TotalCost != oldAperture.TotalCost {
		updates["total_cost"] = e.TotalCost
	}
	if e.IsTotalCostCustom != oldAperture.IsTotalCostCustom {
		updates["is_total_cost_custom"] = e.IsTotalCostCustom
	}
	if e.TotalProfit != oldAperture.TotalProfit {
		updates["total_profit"] = e.TotalProfit
	}
	if e.IsTotalProfitCustom != oldAperture.IsTotalProfitCustom {
		updates["is_total_profit_custom"] = e.IsTotalProfitCustom
	}
	if e.QuoteID != oldAperture.QuoteID {
		updates["quote_id"] = e.QuoteID
	}
	if e.ProfileID != oldAperture.ProfileID {
		updates["profile_id"] = e.ProfileID
	}
	if e.IsProfileIdCustom != oldAperture.IsProfileIdCustom {
		updates["is_profile_id_custom"] = e.IsProfileIdCustom
	}
	if e.GlassID != oldAperture.GlassID {
		updates["glass_id"] = e.GlassID
	}
	if e.IsGlassIdCustom != oldAperture.IsGlassIdCustom {
		updates["is_glass_id_custom"] = e.IsGlassIdCustom
	}
	if e.ColorID != oldAperture.ColorID {
		updates["color_id"] = e.ColorID
	}
	if e.IsColorIdCustom != oldAperture.IsColorIdCustom {
		updates["is_color_id_custom"] = e.IsColorIdCustom
	}
	if e.ParentApertureID != oldAperture.ParentApertureID {
		updates["parent_aperture_id"] = e.ParentApertureID
	}

	return updates
}
