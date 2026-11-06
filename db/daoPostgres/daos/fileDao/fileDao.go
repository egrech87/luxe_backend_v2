package fileDao

import (
	"Quote/db/daoPostgres/baseDao"
	"Quote/db/daoPostgres/daoQuery"
	"Quote/db/entity/baseEntity"
	"Quote/db/entity/entities/fileEntity"
	"context"
	"time"

	"gorm.io/gorm"
)

type FileDao struct {
	baseDao.BaseDao
}

var instance *FileDao

func GetInstance() *FileDao {
	if instance == nil {
		instance = new(FileDao)
	}

	return instance
}

func (d *FileDao) TableName() string {
	return d.NewEntity().TableName()
}
func (d *FileDao) NewEntity() baseEntity.BaseEntityI {
	return fileEntity.New()
}

func (d *FileDao) ApplySelect(tx *gorm.DB, subItemList []string, excludeFields []string) (*gorm.DB, error) {
	fieldNames := d.FieldNamesForTransaction(d, nil)
	tx.Select(fieldNames)

	return tx, nil
}
func (d *FileDao) ApplyFilters(tx *gorm.DB, filters []daoQuery.Filter) (*gorm.DB, error) {
	return tx, nil
}

// #region Basic CRUD operations

func (d *FileDao) Create(ctx context.Context, file baseEntity.BaseEntityI, createdBy uint64) error {
	return d.BaseDao.Create(ctx, file, createdBy)
}

func (d *FileDao) Update(ctx context.Context, id uint64, file baseEntity.BaseEntityI, fieldsToUpdate []string, updatedBy uint64) error {
	return d.BaseDao.Update(ctx, d.TableName(), id, file, fieldsToUpdate, updatedBy)
}

func (d *FileDao) SoftDelete(ctx context.Context, id uint64, deletedBy uint64) (deletedAt time.Time, err error) {
	return d.BaseDao.SoftDelete(ctx, id, fileEntity.TableName(), deletedBy)
}

func (d *FileDao) ReadByID(ctx context.Context, id uint64, selector *daoQuery.Selector) (item baseEntity.BaseEntityI, found bool, err error) {
	return d.BaseDao.ReadByID(ctx, d, id, selector)
}

func (d *FileDao) Read(ctx context.Context, pagination *daoQuery.Pagination, selector *daoQuery.Selector) error {
	return d.ReadMultiple(d.GetBasicTransaction(ctx, d.NewEntity, false), d, pagination, selector)
}

// #endregion

// #region Scopes

// #region WHERE clause

//#endregion

//#endregion
