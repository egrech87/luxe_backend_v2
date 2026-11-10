package apertureDao

import (
	"Quote/db"
	"Quote/db/daoPostgres/baseDao"
	"Quote/db/daoPostgres/daoQuery"
	"Quote/db/daoPostgres/daoQuery/enums/daoOperator"
	"Quote/db/daoPostgres/daos/colorDao"
	"Quote/db/daoPostgres/daos/glassDao"
	"Quote/db/daoPostgres/daos/profileDao"
	"Quote/db/daoPostgres/daos/quoteDao"
	"Quote/db/daoPostgres/daos/userDao"
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/apertureEntity"
	"Quote/utils/anyUtils"
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ApertureDao struct {
	baseDao.BaseDao
}

var instance *ApertureDao

func GetInstance() *ApertureDao {
	if instance == nil {
		instance = new(ApertureDao)
	}

	return instance
}

func (d *ApertureDao) TableName() string {
	return d.NewEntity().TableName()
}
func (d *ApertureDao) NewEntity() baseEntity.BaseEntityI {
	return apertureEntity.New()
}

func (d *ApertureDao) ApplySelect(tx *gorm.DB, subItemList []string, excludeFields []string) (*gorm.DB, error) {
	fieldNames := d.FieldNamesForTransaction(d, nil)
	tableName := d.TableName()

	for _, subItem := range subItemList {
		switch subItem {
		case "quotes":
			quoteDao := quoteDao.GetInstance()
			quoteTableName := quoteDao.TableName()
			quoteFieldNames := quoteDao.FieldNamesForTransaction(quoteDao, nil)
			fieldNames = append(fieldNames, quoteFieldNames...)
			tx.Joins(fmt.Sprintf("left join %s on %s.quote_id = %s.id", quoteTableName, tableName, quoteTableName))
		case "profiles":
			profileDao := profileDao.GetInstance()
			profileTableName := profileDao.TableName()
			profileFieldNames := profileDao.FieldNamesForTransaction(profileDao, nil)
			fieldNames = append(fieldNames, profileFieldNames...)
			tx.Joins(fmt.Sprintf("left join %s on aperture_profile_id(%s.id) = %s.id", profileTableName, tableName, profileTableName))
		case "glasses":
			glassDao := glassDao.GetInstance()
			glassTableName := glassDao.TableName()
			glassFieldNames := glassDao.FieldNamesForTransaction(glassDao, nil)
			fieldNames = append(fieldNames, glassFieldNames...)
			tx.Joins(fmt.Sprintf("left join %s on aperture_glass_id(%s.id) = %s.id", glassTableName, tableName, glassTableName))
		case "colors":
			colorDao := colorDao.GetInstance()
			colorTableName := colorDao.TableName()
			colorFieldNames := colorDao.FieldNamesForTransaction(colorDao, nil)
			fieldNames = append(fieldNames, colorFieldNames...)
			tx.Joins(fmt.Sprintf("left join %s on aperture_color_id(%s.id) = %s.id", colorTableName, tableName, colorTableName))
		case "createuser":
			userDao := userDao.GetInstance()
			userTableName := userDao.TableName()
			userTableAlias := subItem
			userFieldNames := userDao.FieldNamesForTransaction(userDao, &userTableAlias)
			fieldNames = append(fieldNames, userFieldNames...)
			tx.Joins(fmt.Sprintf("left join %s as %s on %s.created_by = %s.id", userTableName, userTableAlias, tableName, userTableAlias))
		case "updateuser":
			userDao := userDao.GetInstance()
			userTableName := userDao.TableName()
			userTableAlias := subItem
			userFieldNames := userDao.FieldNamesForTransaction(userDao, &userTableAlias)
			fieldNames = append(fieldNames, userFieldNames...)
			tx.Joins(fmt.Sprintf("left join %s as %s on aperture_updated_by(%s.id) = %s.id", userTableName, userTableAlias, tableName, userTableAlias))
		}
	}

	// Add functions
	if !slices.Contains(excludeFields, "friendly_name") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_friendly_name(apertures.id) as %s_friendly_name", tableName))
	}
	if !slices.Contains(excludeFields, "image_id") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_image_id(apertures.id) as %s_image_id", tableName))
	}
	if !slices.Contains(excludeFields, "type") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_type(apertures.id) as %s_type", tableName))
	}
	if !slices.Contains(excludeFields, "handle") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_handle(apertures.id) as %s_handle", tableName))
	}
	if !slices.Contains(excludeFields, "area") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_area(apertures.id) as %s_area", tableName))
	}
	if !slices.Contains(excludeFields, "overall_u_value") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_overall_u_value(apertures.id) as %s_overall_u_value", tableName))
	}
	if !slices.Contains(excludeFields, "profile_id") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_profile_id(apertures.id) as %s_profile_id", tableName))
	}
	if !slices.Contains(excludeFields, "color_price_factor") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_color_price_factor(apertures.id) as %s_color_price_factor", tableName))
	}
	if !slices.Contains(excludeFields, "glass_price_factor") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_glass_price_factor(apertures.id) as %s_glass_price_factor", tableName))
	}
	if !slices.Contains(excludeFields, "color_id") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_color_id(apertures.id) as %s_color_id", tableName))
	}
	if !slices.Contains(excludeFields, "glass_id") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_glass_id(apertures.id) as %s_glass_id", tableName))
	}
	if !slices.Contains(excludeFields, "matrix_cost") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_matrix_cost(apertures.id) as %s_matrix_cost", tableName))
	}
	if !slices.Contains(excludeFields, "nets_quantity") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_nets_quantity(apertures.id) as %s_nets_quantity", tableName))
	}
	if !slices.Contains(excludeFields, "nets_price") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_nets_price(apertures.id) as %s_nets_price", tableName))
	}
	if !slices.Contains(excludeFields, "nets_cost") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_nets_cost(apertures.id) as %s_nets_cost", tableName))
	}
	if !slices.Contains(excludeFields, "color_cost") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_color_cost(apertures.id) as %s_color_cost", tableName))
	}
	if !slices.Contains(excludeFields, "glass_cost") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_glass_cost(apertures.id) as %s_glass_cost", tableName))
	}
	if !slices.Contains(excludeFields, "other_cost") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_other_cost(apertures.id) as %s_other_cost", tableName))
	}
	if !slices.Contains(excludeFields, "other_cost_comment") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_other_cost_comment(apertures.id) as %s_other_cost_comment", tableName))
	}
	if !slices.Contains(excludeFields, "supplier_cost") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_supplier_cost(apertures.id) as %s_supplier_cost", tableName))
	}
	if !slices.Contains(excludeFields, "shipping_cost") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_shipping_cost(apertures.id) as %s_shipping_cost", tableName))
	}
	if !slices.Contains(excludeFields, "installation_per_unit") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_installation_per_unit(apertures.id) as %s_installation_per_unit", tableName))
	}
	if !slices.Contains(excludeFields, "installation_cost") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_installation_cost(apertures.id) as %s_installation_cost", tableName))
	}
	if !slices.Contains(excludeFields, "contingency_cost") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_contingency_cost(apertures.id) as %s_contingency_cost", tableName))
	}
	if !slices.Contains(excludeFields, "item_cost") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_item_cost(apertures.id) as %s_item_cost", tableName))
	}
	if !slices.Contains(excludeFields, "profit_margin") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_profit_margin(apertures.id) as %s_profit_margin", tableName))
	}
	if !slices.Contains(excludeFields, "selling_price") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_selling_price(apertures.id) as %s_selling_price", tableName))
	}
	if !slices.Contains(excludeFields, "discounted_selling_price") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_discounted_selling_price(apertures.id) as %s_discounted_selling_price", tableName))
	}
	if !slices.Contains(excludeFields, "total_area") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_total_area(apertures.id) as %s_total_area", tableName))
	}
	if !slices.Contains(excludeFields, "total_selling_price") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_total_selling_price(apertures.id) as %s_total_selling_price", tableName))
	}
	if !slices.Contains(excludeFields, "total_discounted_selling_price") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_total_discounted_selling_price(apertures.id) as %s_total_discounted_selling_price", tableName))
	}
	if !slices.Contains(excludeFields, "total_cost") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_total_cost(apertures.id) as %s_total_cost", tableName))
	}
	if !slices.Contains(excludeFields, "total_profit") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_total_profit(apertures.id) as %s_total_profit", tableName))
	}
	if !slices.Contains(excludeFields, "updated_at") {
		fieldNames = append(fieldNames, fmt.Sprintf("aperture_updated_at(apertures.id) as %s_updated_at", tableName))
	}
	if !slices.Contains(excludeFields, "is_combined") {
		fieldNames = append(fieldNames, fmt.Sprintf("is_aperture_combined(apertures.id) as %s_is_combined", tableName))
	}

	tx.Select(fieldNames)

	return tx, nil
}
func (d *ApertureDao) ApplyFilters(tx *gorm.DB, filters []daoQuery.Filter) (*gorm.DB, error) {
	for _, filter := range filters {
		switch filter.Field {
		case "friendly_name":
			tx.Scopes(d.ByFriendlyName(filter.Operator, filter.Value.(string)))
		case "quote_id":
			value, err := anyUtils.ToUint64(filter.Value)
			if err != nil {
				return tx, err
			}

			tx.Scopes(d.ByQuoteId(filter.Operator, value))
		case "parent_aperture_id":
			switch filter.Operator {
			case daoOperator.Is:
				tx.Scopes(d.ByParentApertureId(filter.Operator, nil))
			default:
				value, err := anyUtils.ToUInt64Ptr(filter.Value)
				if err != nil {
					return tx, err
				}

				tx.Scopes(d.ByParentApertureId(filter.Operator, value))
			}
		}
	}
	return tx, nil
}

// #region Basic CRUD operations

func (d *ApertureDao) Create(ctx context.Context, aperture baseEntity.BaseEntityI, createdBy uint64) error {
	return d.BaseDao.Create(ctx, aperture, createdBy)
}

func (d *ApertureDao) Update(ctx context.Context, id uint64, aperture baseEntity.BaseEntityI, fieldsToUpdate []string, updatedBy uint64) error {
	return d.BaseDao.Update(ctx, d.TableName(), id, aperture, fieldsToUpdate, updatedBy)
}

func (d *ApertureDao) SoftDelete(ctx context.Context, id uint64, deletedBy uint64) (deletedAt time.Time, err error) {
	deletedAt = time.Now()

	softDeletion := map[string]any{}
	softDeletion["updated_by"] = deletedBy
	softDeletion["updated_at"] = &deletedAt
	softDeletion["deleted_by"] = deletedBy
	softDeletion["deleted_at"] = &deletedAt

	tableName := d.TableName()

	err = db.GetDB().WithContext(ctx).Table(tableName).Where(
		db.GetDB().Where(fmt.Sprintf("%s.id = ?", tableName), id).Or(fmt.Sprintf("%s.parent_aperture_id = ?", tableName), id),
	).Scopes(d.NonDeleted(tableName)).Updates(softDeletion).Error

	return deletedAt, err
}

func (d *ApertureDao) ReadByID(ctx context.Context, id uint64, selector *daoQuery.Selector) (item baseEntity.BaseEntityI, found bool, err error) {
	return d.BaseDao.ReadByID(ctx, d, id, selector)
}

func (d *ApertureDao) Read(ctx context.Context, pagination *daoQuery.Pagination, selector *daoQuery.Selector) error {
	return d.ReadMultiple(d.GetBasicTransaction(ctx, d.NewEntity, false), d, pagination, selector)
}

// #endregion

// Duplicate creates `amount` copies of the aperture with the given id, each along with copies of all its non-deleted child apertures.
// The whole operation runs in a single transaction, so either every copy is created or none is.
// It returns the id of the first copy created, and found=false if the source does not exist.
func (d *ApertureDao) Duplicate(ctx context.Context, id uint64, amount int, createdBy uint64) (firstNewId uint64, found bool, err error) {
	tableName := d.TableName()
	found = true

	err = db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		source := apertureEntity.New()
		result := tx.Table(tableName).Scopes(d.ById(tableName, id), d.NonDeleted(tableName)).Limit(1).Find(source)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			found = false
			return nil
		}

		children := []*apertureEntity.ApertureEntity{}
		if err := tx.Table(tableName).Scopes(d.ByParentApertureId(daoOperator.Equals, &id), d.NonDeleted(tableName)).Order(tableName + ".id Asc").Find(&children).Error; err != nil {
			return err
		}

		prepare := func(item *apertureEntity.ApertureEntity) {
			item.SetID(0)
			item.SetCreatedAt(time.Time{})
			item.SetUpdatedAt(time.Time{})
			item.SetDeletedAt(nil)
			item.SetDeletedBy(0)
			item.SetCreatedBy(createdBy)
			item.SetUpdatedBy(createdBy)
		}

		for i := 0; i < amount; i++ {
			// Work on copies so that the source stays untouched between iterations
			newAperture := *source
			prepare(&newAperture)
			if err := tx.Create(&newAperture).Error; err != nil {
				return err
			}
			newId := newAperture.GetID()
			if i == 0 {
				firstNewId = newId
			}

			for _, child := range children {
				newChild := *child
				prepare(&newChild)
				newChild.SetParentApertureID(&newId)
				if err := tx.Create(&newChild).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	return firstNewId, found, err
}

func (d *ApertureDao) SoftDeleteByQuoteId(ctx context.Context, quoteId uint64, deletedAt time.Time, deletedBy uint64) error {
	softDeletion := map[string]any{}
	softDeletion["updated_by"] = deletedBy
	softDeletion["updated_at"] = &deletedAt
	softDeletion["deleted_by"] = deletedBy
	softDeletion["deleted_at"] = &deletedAt

	tableName := d.TableName()

	return db.GetDB().WithContext(ctx).Table(tableName).Scopes(d.ByQuoteId(daoOperator.Equals, quoteId), d.NonDeleted(tableName)).Updates(softDeletion).Error
}

// Fetches the main apertures only for a particular quote
func (d *ApertureDao) ReadMainByQuoteId(ctx context.Context, quoteId uint64, pagination *daoQuery.Pagination, selector *daoQuery.Selector) error {
	return d.ReadMultiple(d.GetBasicTransaction(ctx, d.NewEntity, false).Scopes(d.ByQuoteId(daoOperator.Equals, quoteId), d.ByParentApertureId(daoOperator.Is, nil)).Order("block Asc, unit_name Asc, unit_number Asc"), d, pagination, selector)
}

// Fetches the child apertures only for a particular quote
func (d *ApertureDao) ReadChildrenByQuoteId(ctx context.Context, quoteId uint64, pagination *daoQuery.Pagination, selector *daoQuery.Selector) error {
	return d.ReadMultiple(d.GetBasicTransaction(ctx, d.NewEntity, false).Scopes(d.ByQuoteId(daoOperator.Equals, quoteId), d.ByParentApertureId(daoOperator.IsNot, nil)).Order("apertures.created_at Asc"), d, pagination, selector)
}

// Fetches all the apertures for a particular quote
func (d *ApertureDao) ReadByQuoteId(ctx context.Context, quoteId uint64, pagination *daoQuery.Pagination, selector *daoQuery.Selector) error {
	return d.ReadMultiple(d.GetBasicTransaction(ctx, d.NewEntity, false).Scopes(d.ByQuoteId(daoOperator.Equals, quoteId)).Order("block Asc, unit_name Asc, unit_number Asc"), d, pagination, selector)
}

// #region Scopes

// #region WHERE clause

func (d *ApertureDao) ByFriendlyName(op daoOperator.Operator, friendlyName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("UPPER(%s.friendly_name) LIKE UPPER('%%%s%%')", d.TableName(), strings.TrimSpace(friendlyName)))
	}
}
func (d *ApertureDao) ByDatasheetId(op daoOperator.Operator, datasheet *uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		profileTableName := profileDao.GetInstance().TableName()

		if op == daoOperator.IsNot {
			return db.Where(fmt.Sprintf("%s.datasheet_id IS NOT NULL", profileTableName))
		} else if datasheet != nil {
			return db.Where(fmt.Sprintf("%s.datasheet_id = %d", profileTableName, *datasheet))
		} else {
			return db.Where(fmt.Sprintf("%s.datasheet_id IS NULL", profileTableName))
		}
	}
}
func (d *ApertureDao) ByQuoteId(op daoOperator.Operator, quoteId uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s.quote_id = %d", d.TableName(), quoteId))
	}
}
func (d *ApertureDao) ByProfileId(op daoOperator.Operator, profileId *uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if op == daoOperator.IsNot {
			return db.Where(fmt.Sprintf("%s.profile_id IS NOT NULL", d.TableName()))
		} else if profileId != nil {
			return db.Where(fmt.Sprintf("%s.profile_id = %d", d.TableName(), *profileId))
		} else {
			return db.Where(fmt.Sprintf("%s.profile_id IS NULL", d.TableName()))
		}
	}
}
func (d *ApertureDao) ByParentApertureId(op daoOperator.Operator, parentApertureId *uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if parentApertureId != nil {
			return db.Where(fmt.Sprintf("%s.parent_aperture_id = %d", d.TableName(), *parentApertureId))
		} else {
			if op == daoOperator.Is {
				return db.Where(fmt.Sprintf("%s.parent_aperture_id IS NULL", d.TableName()))
			} else {
				return db.Where(fmt.Sprintf("%s.parent_aperture_id IS NOT NULL", d.TableName()))
			}
		}
	}
}

//#endregion

//#endregion
