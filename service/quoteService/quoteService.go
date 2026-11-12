package quoteService

import (
	"Quote/db"
	"Quote/db/daoPostgres/daos/apertureDao"
	"Quote/db/daoPostgres/daos/quoteDao"
	"Quote/db/entity/entities/apertureEntity"
	"Quote/db/entity/entities/quoteEntity"
	"Quote/utils/anyUtils"
	"context"
	"time"

	"gorm.io/gorm"
)

// Duplicate creates a copy of the quote with the given id, together with copies of all its non-deleted apertures
// and their child apertures. The copy gets the next quote number.
// Everything happens in a single transaction, so either the whole quote is copied or nothing is.
// It returns the id of the new quote, and found=false if the source quote does not exist.
func Duplicate(ctx context.Context, quoteId uint64, createdBy uint64) (newQuoteId uint64, found bool, err error) {
	quoteTable := quoteDao.GetInstance().TableName()
	apertureTable := apertureDao.GetInstance().TableName()
	found = true

	err = db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		quote := quoteEntity.New()
		result := tx.Table(quoteTable).Where(quoteTable+".id = ?", quoteId).Where(quoteTable + ".deleted_at IS NULL").Limit(1).Find(quote)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			found = false
			return nil
		}

		apertures := []*apertureEntity.ApertureEntity{}
		if err := tx.Table(apertureTable).Where(apertureTable+".quote_id = ?", quoteId).Where(apertureTable + ".deleted_at IS NULL").Order(apertureTable + ".id ASC").Find(&apertures).Error; err != nil {
			return err
		}

		// The copy gets the next quote number
		nextNumber := map[string]any{}
		if err := tx.Raw("Select quote_next_number() as next_number").Find(&nextNumber).Error; err != nil {
			return err
		}

		quote.SetID(0)
		quote.SetCreatedAt(time.Time{})
		quote.SetUpdatedAt(time.Time{})
		quote.SetDeletedAt(nil)
		quote.SetDeletedBy(0)
		quote.SetCreatedBy(createdBy)
		quote.SetUpdatedBy(createdBy)
		number, err := anyUtils.ToInt(nextNumber["next_number"])
		if err != nil {
			return err
		}
		quote.Number = number

		if err := tx.Create(quote).Error; err != nil {
			return err
		}
		newQuoteId = quote.GetID()

		prepare := func(aperture *apertureEntity.ApertureEntity) {
			aperture.SetID(0)
			aperture.SetCreatedAt(time.Time{})
			aperture.SetUpdatedAt(time.Time{})
			aperture.SetDeletedAt(nil)
			aperture.SetDeletedBy(0)
			aperture.SetCreatedBy(createdBy)
			aperture.SetUpdatedBy(createdBy)
			aperture.SetQuoteID(&newQuoteId)
		}

		// Copy the main apertures first, remembering which new aperture each old one became
		newIdByOldId := map[uint64]uint64{}
		for _, aperture := range apertures {
			if aperture.GetParentApertureID() != nil {
				continue
			}

			oldId := aperture.GetID()
			prepare(aperture)
			if err := tx.Create(aperture).Error; err != nil {
				return err
			}
			newIdByOldId[oldId] = aperture.GetID()
		}

		// Then the child apertures, attached to the copy of their parent
		for _, aperture := range apertures {
			if aperture.GetParentApertureID() == nil {
				continue
			}

			newParentId, ok := newIdByOldId[*aperture.GetParentApertureID()]
			if !ok {
				continue // The parent is not part of the quote (e.g. deleted), so the child has nothing to attach to
			}

			prepare(aperture)
			aperture.SetParentApertureID(&newParentId)
			if err := tx.Create(aperture).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return newQuoteId, found, err
}
