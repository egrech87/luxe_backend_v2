package glassDao

import (
	"Quote/db/daoPostgres/baseDao"
	"Quote/db/daoPostgres/daoQuery"
	"Quote/db/daoPostgres/daoQuery/enums/daoOperator"
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/glassEntity"
	"Quote/utils/anyUtils"
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type GlassDao struct {
	baseDao.BaseDao
}

var instance *GlassDao

func GetInstance() *GlassDao {
	if instance == nil {
		instance = new(GlassDao)
	}

	return instance
}

func (d *GlassDao) TableName() string {
	return d.NewEntity().TableName()
}
func (d *GlassDao) NewEntity() baseEntity.BaseEntityI {
	return glassEntity.New()
}

func (d *GlassDao) ApplySelect(tx *gorm.DB, subItemList []string, excludeFields []string) (*gorm.DB, error) {
	fieldNames := d.FieldNamesForTransaction(d, nil)
	tx.Select(fieldNames)

	return tx, nil
}
func (d *GlassDao) ApplyFilters(tx *gorm.DB, filters []daoQuery.Filter) (*gorm.DB, error) {
	for _, filter := range filters {
		switch filter.Field {
		case "friendly_name":
			tx.Scopes(d.ByFriendlyName(filter.Operator, filter.Value.(string)))
		case "type":
			value, err := anyUtils.ToIntPtr(filter.Value)
			if err != nil {
				return tx, err
			}

			tx.Scopes(d.ByType(filter.Operator, value))
		}
	}
	return tx, nil
}

// #region Basic CRUD operations

func (d *GlassDao) Create(ctx context.Context, glass baseEntity.BaseEntityI, createdBy uint64) error {
	return d.BaseDao.Create(ctx, glass, createdBy)
}

func (d *GlassDao) Duplicate(ctx context.Context, id uint64, createdBy uint64) (uint64, bool, error) {
	return d.BaseDao.Duplicate(ctx, d, id, createdBy)
}

func (d *GlassDao) Update(ctx context.Context, id uint64, glass baseEntity.BaseEntityI, fieldsToUpdate []string, updatedBy uint64) error {
	return d.BaseDao.Update(ctx, d.TableName(), id, glass, fieldsToUpdate, updatedBy)
}

func (d *GlassDao) SoftDelete(ctx context.Context, id uint64, deletedBy uint64) (deletedAt time.Time, err error) {
	return d.BaseDao.SoftDelete(ctx, id, glassEntity.TableName(), deletedBy)
}

func (d *GlassDao) ReadByID(ctx context.Context, id uint64, selector *daoQuery.Selector) (item baseEntity.BaseEntityI, found bool, err error) {
	return d.BaseDao.ReadByID(ctx, d, id, selector)
}

func (d *GlassDao) Read(ctx context.Context, pagination *daoQuery.Pagination, selector *daoQuery.Selector) error {
	return d.ReadMultiple(d.GetBasicTransaction(ctx, d.NewEntity, false), d, pagination, selector)
}

// #endregion

// #region Scopes

// #region WHERE clause

func (*GlassDao) ByFriendlyName(op daoOperator.Operator, friendlyName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("UPPER(friendly_name) LIKE UPPER('%%%s%%')", strings.TrimSpace(friendlyName)))
	}
}
func (*GlassDao) ByType(op daoOperator.Operator, glassType *int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if glassType != nil && *glassType >= 0 {
			return db.Where(fmt.Sprintf("type = %d", *glassType))
		} else {
			return db.Where("type IS NULL")
		}
	}
}

//#endregion

//#endregion
