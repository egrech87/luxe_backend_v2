package colorDao

import (
	"Quote/db/daoPostgres/baseDao"
	"Quote/db/daoPostgres/daoQuery"
	"Quote/db/daoPostgres/daoQuery/enums/daoOperator"
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/colorEntity"
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ColorDao struct {
	baseDao.BaseDao
}

var instance *ColorDao

func GetInstance() *ColorDao {
	if instance == nil {
		instance = new(ColorDao)
	}

	return instance
}

func (d *ColorDao) TableName() string {
	return d.NewEntity().TableName()
}
func (d *ColorDao) NewEntity() baseEntity.BaseEntityI {
	return colorEntity.New()
}

func (d *ColorDao) ApplySelect(tx *gorm.DB, subItemList []string, excludeFields []string) (*gorm.DB, error) {
	fieldNames := d.FieldNamesForTransaction(d, nil)
	tx.Select(fieldNames)

	return tx, nil
}
func (d *ColorDao) ApplyFilters(tx *gorm.DB, filters []daoQuery.Filter) (*gorm.DB, error) {
	for _, filter := range filters {
		switch filter.Field {
		case "friendly_name":
			tx.Scopes(d.ByFriendlyName(filter.Operator, filter.Value.(string)))
		case "code":
			tx.Scopes(d.ByCode(filter.Operator, filter.Value.(string)))
		}
	}
	return tx, nil
}

// #region Basic CRUD operations

func (d *ColorDao) Create(ctx context.Context, color baseEntity.BaseEntityI, createdBy uint64) error {
	return d.BaseDao.Create(ctx, color, createdBy)
}

func (d *ColorDao) Duplicate(ctx context.Context, id uint64, createdBy uint64) (uint64, bool, error) {
	return d.BaseDao.Duplicate(ctx, d, id, createdBy)
}

func (d *ColorDao) Update(ctx context.Context, id uint64, color baseEntity.BaseEntityI, fieldsToUpdate []string, updatedBy uint64) error {
	return d.BaseDao.Update(ctx, d.TableName(), id, color, fieldsToUpdate, updatedBy)
}

func (d *ColorDao) SoftDelete(ctx context.Context, id uint64, deletedBy uint64) (deletedAt time.Time, err error) {
	return d.BaseDao.SoftDelete(ctx, id, colorEntity.TableName(), deletedBy)
}

func (d *ColorDao) ReadByID(ctx context.Context, id uint64, selector *daoQuery.Selector) (item baseEntity.BaseEntityI, found bool, err error) {
	return d.BaseDao.ReadByID(ctx, d, id, selector)
}

func (d *ColorDao) Read(ctx context.Context, pagination *daoQuery.Pagination, selector *daoQuery.Selector) error {
	return d.ReadMultiple(d.GetBasicTransaction(ctx, d.NewEntity, false), d, pagination, selector)
}

// #endregion

// #region Scopes

// #region WHERE clause

func (*ColorDao) ByFriendlyName(op daoOperator.Operator, friendlyName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("UPPER(friendly_name) LIKE UPPER('%%%s%%')", strings.TrimSpace(friendlyName)))
	}
}
func (*ColorDao) ByCode(op daoOperator.Operator, code string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("UPPER(code) LIKE UPPER('%%%s%%')", strings.TrimSpace(code)))
	}
}

//#endregion

//#endregion
