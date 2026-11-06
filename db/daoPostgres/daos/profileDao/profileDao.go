package profileDao

import (
	"Quote/db/daoPostgres/baseDao"
	"Quote/db/daoPostgres/daoQuery"
	"Quote/db/daoPostgres/daoQuery/enums/daoOperator"
	"Quote/db/daoPostgres/daos/supplierDao"
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/profileEntity"
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ProfileDao struct {
	baseDao.BaseDao
}

var instance *ProfileDao

func GetInstance() *ProfileDao {
	if instance == nil {
		instance = new(ProfileDao)
	}

	return instance
}

func (d *ProfileDao) TableName() string {
	return d.NewEntity().TableName()
}
func (d *ProfileDao) NewEntity() baseEntity.BaseEntityI {
	return profileEntity.New()
}

func (d *ProfileDao) ApplySelect(tx *gorm.DB, subItemList []string, excludeFields []string) (*gorm.DB, error) {
	fieldNames := d.FieldNamesForTransaction(d, nil)
	tableName := d.TableName()

	for _, subItem := range subItemList {
		switch subItem {
		case "suppliers":
			supplierDao := supplierDao.GetInstance()
			supplierTableName := supplierDao.TableName()
			supplierFieldNames := supplierDao.FieldNamesForTransaction(supplierDao, nil)
			fieldNames = append(fieldNames, supplierFieldNames...)
			tx.Joins(fmt.Sprintf("left join %s on %s.supplier_id = %s.id", supplierTableName, tableName, supplierTableName))
		}
	}

	tx.Select(fieldNames)

	return tx, nil
}
func (d *ProfileDao) ApplyFilters(tx *gorm.DB, filters []daoQuery.Filter) (*gorm.DB, error) {
	for _, filter := range filters {
		switch filter.Field {
		case "friendly_name":
			tx.Scopes(d.ByFriendlyName(filter.Operator, filter.Value.(string)))
		}
	}
	return tx, nil
}

// #region Basic CRUD operations

func (d *ProfileDao) Create(ctx context.Context, profile baseEntity.BaseEntityI, createdBy uint64) error {
	return d.BaseDao.Create(ctx, profile, createdBy)
}

func (d *ProfileDao) Update(ctx context.Context, id uint64, profile baseEntity.BaseEntityI, fieldsToUpdate []string, updatedBy uint64) error {
	return d.BaseDao.Update(ctx, d.TableName(), id, profile, fieldsToUpdate, updatedBy)
}

func (d *ProfileDao) SoftDelete(ctx context.Context, id uint64, deletedBy uint64) (deletedAt time.Time, err error) {
	return d.BaseDao.SoftDelete(ctx, id, profileEntity.TableName(), deletedBy)
}

func (d *ProfileDao) ReadByID(ctx context.Context, id uint64, selector *daoQuery.Selector) (item baseEntity.BaseEntityI, found bool, err error) {
	return d.BaseDao.ReadByID(ctx, d, id, selector)
}

func (d *ProfileDao) Read(ctx context.Context, pagination *daoQuery.Pagination, selector *daoQuery.Selector) error {
	return d.ReadMultiple(d.GetBasicTransaction(ctx, d.NewEntity, false), d, pagination, selector)
}

// #endregion

// #region Scopes

// #region WHERE clause

func (d *ProfileDao) ByFriendlyName(op daoOperator.Operator, friendlyName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("UPPER(%s.friendly_name) LIKE UPPER('%%%s%%')", d.TableName(), strings.TrimSpace(friendlyName)))
	}
}

//#endregion

//#endregion
