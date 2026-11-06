package baseDao

import (
	"Quote/db"
	"Quote/db/daoPostgres/daoQuery"
	"Quote/db/entity/baseEntity"
	"fmt"

	"Quote/utils/arrayUtils"
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseDaoI interface {
	TableName() string
	NewEntity() baseEntity.BaseEntityI
	ApplySelect(tx *gorm.DB, subItemList []string, excludeFields []string) (*gorm.DB, error)
	ApplyFilters(tx *gorm.DB, filters []daoQuery.Filter) (*gorm.DB, error)
	Create(ctx context.Context, item baseEntity.BaseEntityI, createdBy uint64) error
	Update(ctx context.Context, id uint64, item baseEntity.BaseEntityI, fieldsToUpdate []string, updatedBy uint64) error
	SoftDelete(ctx context.Context, id uint64, deletedBy uint64) (deletedAt time.Time, err error)
	ReadByID(ctx context.Context, id uint64, selector *daoQuery.Selector) (baseEntity.BaseEntityI, bool, error)
	Read(ctx context.Context, pagination *daoQuery.Pagination, selector *daoQuery.Selector) error
}

type BaseDao struct {
}

// #region Basic CRUD operations

func (d *BaseDao) Create(ctx context.Context, item baseEntity.BaseEntityI, createdBy uint64) error {
	item.SetCreatedBy(createdBy)
	item.SetUpdatedBy(createdBy)

	tx := db.GetDB().WithContext(ctx).Create(item)

	return tx.Error
}
func (d *BaseDao) Update(ctx context.Context, tableName string, id uint64, item baseEntity.BaseEntityI, fieldsToUpdate []string, updatedBy uint64) error {
	item.SetUpdatedBy(updatedBy)
	fieldsToUpdate, _ = arrayUtils.AppendUnique(fieldsToUpdate, "updated_by")

	return db.GetDB().WithContext(ctx).Model(&item).Clauses(clause.Returning{}).Scopes(d.ById(tableName, id), d.NonDeleted(tableName)).Select(fieldsToUpdate).Updates(item).Error
}

func (d *BaseDao) SoftDelete(ctx context.Context, id uint64, tableName string, deletedBy uint64) (deletedAt time.Time, err error) {
	deletedAt = time.Now()

	softDeletion := map[string]any{}
	softDeletion["updated_by"] = deletedBy
	softDeletion["updated_at"] = &deletedAt
	softDeletion["deleted_by"] = deletedBy
	softDeletion["deleted_at"] = &deletedAt

	err = db.GetDB().WithContext(ctx).Table(tableName).Scopes(d.ById(tableName, id), d.NonDeleted(tableName)).Updates(softDeletion).Error
	return deletedAt, err
}

func (d *BaseDao) ReadSingle(tx *gorm.DB, dao BaseDaoI, selector *daoQuery.Selector) (item baseEntity.BaseEntityI, found bool, err error) {
	tx, err = d.ApplySelector(tx, dao, selector)
	if err != nil {
		return nil, false, err
	}

	itemMap := map[string]any{}
	result := tx.Find(&itemMap, db.GetDB())

	item = dao.NewEntity()
	item.FillFromDb(itemMap)

	found = true
	if result.RowsAffected == 0 {
		found = false
	}

	return item, found, result.Error
}

func (d *BaseDao) ReadByID(ctx context.Context, dao BaseDaoI, id uint64, selector *daoQuery.Selector) (item baseEntity.BaseEntityI, found bool, err error) {
	return d.ReadSingle(d.getByIdTx(ctx, dao.TableName(), id), dao, selector)
}

func (d *BaseDao) ReadMultiple(tx *gorm.DB, dao BaseDaoI, pagination *daoQuery.Pagination, selector *daoQuery.Selector) (err error) {
	items := []any{}

	tx, err = d.ApplySelector(tx, dao, selector)
	if err != nil {
		return err
	}

	var totalItems int64
	if err = tx.Count(&totalItems).Error; err != nil {
		return err
	}
	pagination.Meta.TotalItems = int(totalItems)

	itemListMap := []map[string]any{}
	// tx = tx.Debug() // Use this to view the SQL in Debug Console
	if err := tx.Scopes(pagination.Paginate(pagination)).Find(&itemListMap).Error; err != nil {
		return err
	}

	for _, itemMap := range itemListMap {
		item := dao.NewEntity()
		if err = item.FillFromDb(itemMap); err != nil {
			return err
		}
		items = append(items, item)
	}

	pagination.Meta.TotalPages = pagination.CalculateTotalPages()
	pagination.Data = items

	return nil
}

// #endregion

// #region Scopes

// #region WHERE clause

func (d *BaseDao) GetBasicTransaction(ctx context.Context, newItemFunc func() baseEntity.BaseEntityI, containDeleted bool) *gorm.DB {
	newItem := newItemFunc()
	tx := db.GetDB().WithContext(ctx).Model(newItem)

	if !containDeleted {
		tx = tx.Scopes(d.NonDeleted(newItem.TableName()))
	}

	return tx
}

func (d *BaseDao) getByIdTx(ctx context.Context, tableName string, id uint64) *gorm.DB {
	return db.GetDB().WithContext(ctx).Table(tableName).Scopes(d.ById(tableName, id))
}

func (*BaseDao) ById(tableName string, id uint64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s.id = ?", tableName), id)
	}
}
func (*BaseDao) NonDeleted(tableName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s.deleted_at is null", tableName))
	}
}

//#endregion

//#endregion

// Provide nil when no alias is needed
func (d *BaseDao) FieldNamesForTransaction(dao BaseDaoI, tableAlias *string) []string {
	var tableName string

	if tableAlias == nil {
		tableName = dao.TableName()
	} else {
		tableName = *tableAlias
	}

	var fieldNames []string

	fieldNames = []string{}
	for _, fieldName := range dao.NewEntity().FieldNames() {
		fieldNames = append(fieldNames, fmt.Sprintf("%s.%s as %s_%s", tableName, fieldName, tableName, fieldName))
	}

	return fieldNames
}

func (*BaseDao) ApplySelector(tx *gorm.DB, dao BaseDaoI, selector *daoQuery.Selector) (*gorm.DB, error) {
	var err error

	// Apply the Select
	{
		excludeFields := []string{}
		subItems := []string{}
		if selector != nil {
			excludeFields = selector.ExcludeFields
			subItems = selector.SubItems
		}

		tx, err = dao.ApplySelect(tx, subItems, excludeFields)
		if err != nil {
			return tx, err
		}
	}

	// Apply the filters
	if selector != nil {
		tx, err = dao.ApplyFilters(tx, selector.Filters)
		if err != nil {
			return tx, err
		}
	}

	return tx, nil
}
