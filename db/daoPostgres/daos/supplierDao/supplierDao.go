package supplierDao

import (
	"Quote/db/daoPostgres/baseDao"
	"Quote/db/daoPostgres/daoQuery"
	"Quote/db/daoPostgres/daoQuery/enums/daoOperator"
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/supplierEntity"
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type SupplierDao struct {
	baseDao.BaseDao
}

var instance *SupplierDao

func GetInstance() *SupplierDao {
	if instance == nil {
		instance = new(SupplierDao)
	}

	return instance
}

func (d *SupplierDao) TableName() string {
	return d.NewEntity().TableName()
}
func (d *SupplierDao) NewEntity() baseEntity.BaseEntityI {
	return supplierEntity.New()
}

func (d *SupplierDao) ApplySelect(tx *gorm.DB, subItemList []string, excludeFields []string) (*gorm.DB, error) {
	fieldNames := d.FieldNamesForTransaction(d, nil)
	tx.Select(fieldNames)

	return tx, nil
}
func (d *SupplierDao) ApplyFilters(tx *gorm.DB, filters []daoQuery.Filter) (*gorm.DB, error) {
	for _, filter := range filters {
		switch filter.Field {
		case "friendly_name":
			tx.Scopes(d.ByFriendlyName(filter.Operator, filter.Value.(string)))
		case "name":
			tx.Scopes(d.ByName(filter.Operator, filter.Value.(string)))
		case "address":
			tx.Scopes(d.ByAddress(filter.Operator, filter.Value.(string)))
		case "email":
			tx.Scopes(d.ByEmail(filter.Operator, filter.Value.(string)))
		case "contact_number":
			tx.Scopes(d.ByContactNumber(filter.Operator, filter.Value.(string)))
		}
	}
	return tx, nil
}

// #region Basic CRUD operations

func (d *SupplierDao) Create(ctx context.Context, supplier baseEntity.BaseEntityI, createdBy uint64) error {
	return d.BaseDao.Create(ctx, supplier, createdBy)
}

func (d *SupplierDao) Update(ctx context.Context, id uint64, supplier baseEntity.BaseEntityI, fieldsToUpdate []string, updatedBy uint64) error {
	return d.BaseDao.Update(ctx, d.TableName(), id, supplier, fieldsToUpdate, updatedBy)
}

func (d *SupplierDao) SoftDelete(ctx context.Context, id uint64, deletedBy uint64) (deletedAt time.Time, err error) {
	return d.BaseDao.SoftDelete(ctx, id, supplierEntity.TableName(), deletedBy)
}

func (d *SupplierDao) ReadByID(ctx context.Context, id uint64, selector *daoQuery.Selector) (item baseEntity.BaseEntityI, found bool, err error) {
	return d.BaseDao.ReadByID(ctx, d, id, selector)
}

func (d *SupplierDao) Read(ctx context.Context, pagination *daoQuery.Pagination, selector *daoQuery.Selector) error {
	return d.ReadMultiple(d.GetBasicTransaction(ctx, d.NewEntity, false), d, pagination, selector)
}

// #endregion

// #region Scopes

// #region WHERE clause

func (*SupplierDao) ByFriendlyName(op daoOperator.Operator, friendlyName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("UPPER(friendly_name) LIKE UPPER('%%%s%%')", strings.TrimSpace(friendlyName)))
	}
}
func (*SupplierDao) ByName(op daoOperator.Operator, name string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("UPPER(name) LIKE UPPER('%%%s%%')", strings.TrimSpace(name)))
	}
}
func (*SupplierDao) ByAddress(op daoOperator.Operator, address string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("UPPER(address) LIKE UPPER('%%%s%%')", strings.TrimSpace(address)))
	}
}
func (*SupplierDao) ByEmail(op daoOperator.Operator, email string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("UPPER(email) LIKE UPPER('%%%s%%')", strings.TrimSpace(email)))
	}
}
func (*SupplierDao) ByContactNumber(op daoOperator.Operator, contactNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("UPPER(contact_number) LIKE UPPER('%%%s%%')", strings.TrimSpace(contactNumber)))
	}
}

//#endregion

//#endregion
