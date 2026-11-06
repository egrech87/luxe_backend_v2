package clientDao

import (
	"Quote/db/daoPostgres/baseDao"
	"Quote/db/daoPostgres/daoQuery"
	"Quote/db/daoPostgres/daoQuery/enums/daoOperator"
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/clientEntity"
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type ClientDao struct {
	baseDao.BaseDao
}

var instance *ClientDao

func GetInstance() *ClientDao {
	if instance == nil {
		instance = new(ClientDao)
	}

	return instance
}

func (d *ClientDao) TableName() string {
	return d.NewEntity().TableName()
}
func (d *ClientDao) NewEntity() baseEntity.BaseEntityI {
	return clientEntity.New()
}

func (d *ClientDao) ApplySelect(tx *gorm.DB, subItemList []string, excludeFields []string) (*gorm.DB, error) {
	fieldNames := d.FieldNamesForTransaction(d, nil)
	tx.Select(fieldNames)

	return tx, nil
}
func (d *ClientDao) ApplyFilters(tx *gorm.DB, filters []daoQuery.Filter) (*gorm.DB, error) {
	for _, filter := range filters {
		switch filter.Field {
		case "name":
			tx.Scopes(d.ByName(filter.Operator, filter.Value.(string)))
		case "surname":
			tx.Scopes(d.BySurname(filter.Operator, filter.Value.(string)))
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

func (d *ClientDao) Create(ctx context.Context, client baseEntity.BaseEntityI, createdBy uint64) error {
	return d.BaseDao.Create(ctx, client, createdBy)
}

func (d *ClientDao) Update(ctx context.Context, id uint64, client baseEntity.BaseEntityI, fieldsToUpdate []string, updatedBy uint64) error {
	return d.BaseDao.Update(ctx, d.TableName(), id, client, fieldsToUpdate, updatedBy)
}

func (d *ClientDao) SoftDelete(ctx context.Context, id uint64, deletedBy uint64) (deletedAt time.Time, err error) {
	return d.BaseDao.SoftDelete(ctx, id, clientEntity.TableName(), deletedBy)
}

func (d *ClientDao) ReadByID(ctx context.Context, id uint64, selector *daoQuery.Selector) (item baseEntity.BaseEntityI, found bool, err error) {
	return d.BaseDao.ReadByID(ctx, d, id, selector)
}

func (d *ClientDao) Read(ctx context.Context, pagination *daoQuery.Pagination, selector *daoQuery.Selector) error {
	return d.ReadMultiple(d.GetBasicTransaction(ctx, d.NewEntity, false), d, pagination, selector)
}

// #endregion

// #region Scopes

// #region WHERE clause

func (*ClientDao) ByName(op daoOperator.Operator, name string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("UPPER(name) LIKE UPPER('%%%s%%')", strings.TrimSpace(name)))
	}
}
func (*ClientDao) BySurname(op daoOperator.Operator, surname string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("UPPER(surname) LIKE UPPER('%%%s%%')", strings.TrimSpace(surname)))
	}
}
func (*ClientDao) ByAddress(op daoOperator.Operator, address string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("UPPER(address) LIKE UPPER('%%%s%%')", strings.TrimSpace(address)))
	}
}
func (*ClientDao) ByEmail(op daoOperator.Operator, email string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("UPPER(email) LIKE UPPER('%%%s%%')", strings.TrimSpace(email)))
	}
}
func (*ClientDao) ByContactNumber(op daoOperator.Operator, contactNumber string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("UPPER(contact_number) LIKE UPPER('%%%s%%')", strings.TrimSpace(contactNumber)))
	}
}

//#endregion

//#endregion
