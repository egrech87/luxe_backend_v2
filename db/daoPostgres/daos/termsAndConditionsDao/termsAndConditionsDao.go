package termsAndConditionsDao

import (
	"Quote/db"
	"Quote/db/entity/entities/termsAndConditionsEntity"
	"context"
	"time"

	"gorm.io/gorm"
)

type TermsAndConditionsDao struct{}

var instance *TermsAndConditionsDao

func GetInstance() *TermsAndConditionsDao {
	if instance == nil {
		instance = new(TermsAndConditionsDao)
	}

	return instance
}

func (d *TermsAndConditionsDao) tableName() string {
	return termsAndConditionsEntity.TableName()
}

// Get returns the current terms and conditions record, if there is one.
func (d *TermsAndConditionsDao) Get(ctx context.Context) (item *termsAndConditionsEntity.TermsAndConditionsEntity, found bool, err error) {
	item = termsAndConditionsEntity.New()

	result := db.GetDB().WithContext(ctx).Table(d.tableName()).Where("deleted_at IS NULL").Order("id DESC").Limit(1).Find(item)
	if result.Error != nil {
		return nil, false, result.Error
	}

	return item, result.RowsAffected > 0, nil
}

// Set makes the given file the current terms and conditions, creating the record if needed.
// It returns the id of the file that was replaced (0 if there was none).
func (d *TermsAndConditionsDao) Set(ctx context.Context, fileId uint64, userId uint64) (replacedFileId uint64, err error) {
	err = db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		current := termsAndConditionsEntity.New()
		result := tx.Table(d.tableName()).Where("deleted_at IS NULL").Order("id DESC").Limit(1).Find(current)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			current = termsAndConditionsEntity.New()
			current.FileID = fileId
			current.SetCreatedBy(userId)
			current.SetUpdatedBy(userId)
			return tx.Create(current).Error
		}

		replacedFileId = current.FileID
		return tx.Table(d.tableName()).Where("id = ?", current.GetID()).Updates(map[string]any{
			"file_id":    fileId,
			"updated_by": userId,
			"updated_at": time.Now(),
		}).Error
	})

	return replacedFileId, err
}

// SoftDelete removes the current terms and conditions record.
// It returns the id of the file it pointed to, and found=false if there was no record.
func (d *TermsAndConditionsDao) SoftDelete(ctx context.Context, deletedBy uint64) (fileId uint64, found bool, err error) {
	err = db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		current := termsAndConditionsEntity.New()
		result := tx.Table(d.tableName()).Where("deleted_at IS NULL").Order("id DESC").Limit(1).Find(current)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return nil
		}

		found = true
		fileId = current.FileID

		now := time.Now()
		return tx.Table(d.tableName()).Where("deleted_at IS NULL").Updates(map[string]any{
			"updated_by": deletedBy,
			"updated_at": now,
			"deleted_by": deletedBy,
			"deleted_at": now,
		}).Error
	})

	return fileId, found, err
}
