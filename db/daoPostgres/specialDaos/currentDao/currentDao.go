package currentDao

import (
	"Quote/db"
	"Quote/db/daoPostgres/baseDao"
	"Quote/db/entity/specialEntities/currentEntity"
	"context"

	"gorm.io/gorm/clause"
)

type CurrentDao struct {
	baseDao.BaseDao
}

var instance *CurrentDao

func GetInstance() *CurrentDao {
	if instance == nil {
		instance = new(CurrentDao)
	}

	return instance
}

// #region Basic CRUD operations

func (d *CurrentDao) Update(ctx context.Context, current *currentEntity.CurrentEntity, fieldsToUpdate []string) error {
	return db.GetDB().WithContext(ctx).Model(&current).Clauses(clause.Returning{}).Select(fieldsToUpdate).Where("1 = 1").Updates(current).Error
}

func (d *CurrentDao) Read(ctx context.Context) (*currentEntity.CurrentEntity, error) {
	item := currentEntity.New()

	result := db.GetDB().WithContext(ctx).Table(currentEntity.TableName()).Take(item)

	return item, result.Error
}

// #endregion
