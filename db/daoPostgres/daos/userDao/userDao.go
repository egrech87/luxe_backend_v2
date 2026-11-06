package userDao

import (
	"Quote/db/daoPostgres/baseDao"
	"Quote/db/daoPostgres/daoQuery"
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/userEntity"
	"context"
	"time"

	"gorm.io/gorm"
)

type UserDao struct {
	baseDao.BaseDao
}

var instance *UserDao

func GetInstance() *UserDao {
	if instance == nil {
		instance = new(UserDao)
	}

	return instance
}

func (d *UserDao) TableName() string {
	return d.NewEntity().TableName()
}
func (d *UserDao) NewEntity() baseEntity.BaseEntityI {
	return userEntity.New()
}

func (d *UserDao) ApplySelect(tx *gorm.DB, subItemList []string, excludeFields []string) (*gorm.DB, error) {
	fieldNames := d.FieldNamesForTransaction(d, nil)
	tx.Select(fieldNames)

	return tx, nil
}
func (d *UserDao) ApplyFilters(tx *gorm.DB, filters []daoQuery.Filter) (*gorm.DB, error) {
	return tx, nil
}

// #region Basic CRUD operations

func (d *UserDao) Create(ctx context.Context, user baseEntity.BaseEntityI, createdBy uint64) error {
	return d.BaseDao.Create(ctx, user, createdBy)
}

func (d *UserDao) Update(ctx context.Context, id uint64, user baseEntity.BaseEntityI, fieldsToUpdate []string, updatedBy uint64) error {
	return d.BaseDao.Update(ctx, d.TableName(), id, user, fieldsToUpdate, updatedBy)
}

func (d *UserDao) SoftDelete(ctx context.Context, id uint64, deletedBy uint64) (deletedAt time.Time, err error) {
	return d.BaseDao.SoftDelete(ctx, id, userEntity.TableName(), deletedBy)
}

func (d *UserDao) ReadByID(ctx context.Context, id uint64, selector *daoQuery.Selector) (item baseEntity.BaseEntityI, found bool, err error) {
	return d.BaseDao.ReadByID(ctx, d, id, selector)
}

func (d *UserDao) Read(ctx context.Context, pagination *daoQuery.Pagination, selector *daoQuery.Selector) error {
	return d.ReadMultiple(d.GetBasicTransaction(ctx, d.NewEntity, false), d, pagination, selector)
}

// #endregion

// #region Scopes

// #region WHERE clause

func (*UserDao) ByFriendlyName(friendlyName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("UPPER(friendly_name) = UPPER(?)", friendlyName)
	}
}

//#endregion

//#endregion
