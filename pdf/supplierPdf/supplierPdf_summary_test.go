package supplierPdf

import (
	"Quote/db/entity/entities/apertureEntity"
	"math"
	"testing"
)

func TestSummarise(t *testing.T) {
	aperture := func(quantity, actualWidth, actualHeight int) *apertureEntity.ApertureEntity {
		a := apertureEntity.New()
		a.Quantity, a.ActualWidth, a.ActualHeight = quantity, actualWidth, actualHeight
		return a
	}

	// The area is the sum of actual width x actual height of each aperture (not multiplied by the quantity),
	// the amount of apertures is the sum of the quantities. Child apertures are not part of the figures.
	main := aperture(2, 1000, 2000)
	main.ChildApertures = []*apertureEntity.ApertureEntity{aperture(5, 400, 500)}

	area, amount := summarise([]*apertureEntity.ApertureEntity{main, aperture(3, 900, 1200), aperture(1, 1500, 1500)})
	if math.Abs(area-5.33) > 1e-9 {
		t.Fatalf("expected a total area of 5.33, got %v", area)
	}
	if amount != 6 {
		t.Fatalf("expected 6 apertures, got %d", amount)
	}

	area, amount = summarise(nil)
	if area != 0 || amount != 0 {
		t.Fatalf("expected no area and no apertures, got %v and %d", area, amount)
	}
}
