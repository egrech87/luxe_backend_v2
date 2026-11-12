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
		if err := tx.Table(apertureTable).Where(apertureTable+".quote_id = ?", quoteId).Where(apertureTable + ".deleted_at IS NULL").Order(apertureTable + ".created_at ASC, " + apertureTable + ".id ASC").Find(&apertures).Error; err != nil {
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

		// Copy the apertures in the order they were created (the list is sorted by creation date, oldest first),
		// so the copies keep the same relative order. A child is attached to the copy of its parent.
		newIdByOldId := map[uint64]uint64{}
		waitingForParent := map[uint64][]*apertureEntity.ApertureEntity{} // children whose parent comes later in the list

		var copyAperture func(aperture *apertureEntity.ApertureEntity, newParentId *uint64) error
		copyAperture = func(aperture *apertureEntity.ApertureEntity, newParentId *uint64) error {
			oldId := aperture.GetID()
			prepare(aperture)
			aperture.SetParentApertureID(newParentId)
			if err := tx.Create(aperture).Error; err != nil {
				return err
			}
			newIdByOldId[oldId] = aperture.GetID()

			// Now that this aperture has a copy, copy any children that were waiting for it
			for _, child := range waitingForParent[oldId] {
				if err := copyAperture(child, &[]uint64{aperture.GetID()}[0]); err != nil {
					return err
				}
			}
			delete(waitingForParent, oldId)

			return nil
		}

		for _, aperture := range apertures {
			parentId := aperture.GetParentApertureID()
			if parentId == nil {
				if err := copyAperture(aperture, nil); err != nil {
					return err
				}
				continue
			}

			if newParentId, ok := newIdByOldId[*parentId]; ok {
				if err := copyAperture(aperture, &newParentId); err != nil {
					return err
				}
			} else {
				waitingForParent[*parentId] = append(waitingForParent[*parentId], aperture)
			}
		}
		// Any child still waiting has a parent that is not part of the quote (e.g. deleted), so it is not copied

		return nil
	})

	return newQuoteId, found, err
}
